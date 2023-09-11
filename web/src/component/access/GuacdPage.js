// Copyright 2023 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import React, {useEffect, useState} from 'react';
import {useLocation} from "react-router-dom";
import {Affix, Button, Dropdown, Menu, message, Modal} from "antd";
import {CloseCircleOutlined, CopyOutlined, ExpandOutlined, WindowsOutlined} from "@ant-design/icons";
import Guacamole from "guacamole-common-js";
import {exitFull, getToken, requestFullScreen, debounce} from "./Util";
import * as Setting from "../../Setting";
import qs from "qs";
import {Base64} from "js-base64";
import Draggable from "react-draggable";
import GuacdClipboard from "./GuacdClipboard";
import './Guacd.css';

let fixedSize = false;

const STATE_IDLE = 0;
const STATE_CONNECTING = 1;
const STATE_WAITING = 2;
const STATE_CONNECTED = 3;
const STATE_DISCONNECTING = 4;
const STATE_DISCONNECTED = 5;

const GuacdPage = () => {
  const location = useLocation();
  let searchParams = new URLSearchParams(location.search);

  let assetOwner = searchParams.get('owner');
  let assetName = searchParams.get('name');
  let protocol = searchParams.get('protocol');
  let width = searchParams.get('width');
  let height = searchParams.get('height');

  if (width && height) {
    fixedSize = true;
  } else {
    width = window.innerWidth;
    height = window.innerHeight;
  }

  let [box, setBox] = useState({width, height});
  let [guacd, setGuacd] = useState({});
  let [session, setSession] = useState({});
  let [clipboardText, setClipboardText] = useState('');
  let [fullScreened, setFullScreened] = useState(false);
  let [clipboardVisible, setClipboardVisible] = useState(false);

  useEffect(() => {
    document.title = assetName;
    createSession();
  }, [assetOwner, assetName]);

  const createSession = async () => {
    const newSession = {
      owner: assetOwner,
      name: assetName,
      mode: 'guacd',
    };

    renderDisplay(assetOwner, assetName, protocol, width, height);
  }

  const renderDisplay = (assetOwner, assetName, protocol, width, height) => {
    let sessionId = "123"
    const wsEndpoint = Setting.ServerUrl.replace("http://", "ws://");
    const wsUrl = `${wsEndpoint}/api/get-asset-tunnel?owner=${assetOwner}&name=${assetName}&`;
    let tunnel = new Guacamole.WebSocketTunnel(wsUrl);
    let client = new Guacamole.Client(tunnel);

    // Handling clipboard content received from a virtual machine.
    client.onclipboard = handleClipboardReceived;

    // Handling client state change events.
    client.onstatechange = (state) => {
      onClientStateChange(state, sessionId);
    };

    client.onerror = onError;
    tunnel.onerror = onError;

    // Get display div from document
    const displayEle = document.getElementById("display");

    // Add client to display div
    const element = client.getDisplay().getElement();
    displayEle.appendChild(element);

    let dpi = 96;
    if (protocol === 'telnet') {
      dpi = dpi * 2;
    }

    let token = getToken();

    let params = {'width': width, 'height': height, 'dpi': dpi, 'X-Auth-Token': token};

    let paramStr = qs.stringify(params);

    client.connect(paramStr);
    let display = client.getDisplay();
    display.onresize = function (width, height) {
      display.scale(Math.min(window.innerHeight / display.getHeight(), window.innerWidth / display.getHeight()));
    }

    const sink = new Guacamole.InputSink();
    displayEle.appendChild(sink.getElement());
    sink.focus();

    const keyboard = new Guacamole.Keyboard(sink.getElement());

    keyboard.onkeydown = (keysym) => {
      client.sendKeyEvent(1, keysym);
      if (keysym === 65288) {
        return false;
      }
    };
    keyboard.onkeyup = (keysym) => {
      client.sendKeyEvent(0, keysym);
    };

    const sinkFocus = debounce(() => {
      sink.focus();
    });

    const mouse = new Guacamole.Mouse(element);

    mouse.onmousedown = mouse.onmouseup = function (mouseState) {
      sinkFocus();
      client.sendMouseState(mouseState);
    }

    mouse.onmousemove = function (mouseState) {
      sinkFocus();
      client.getDisplay().showCursor(false);
      mouseState.x = mouseState.x / display.getScale();
      mouseState.y = mouseState.y / display.getScale();
      client.sendMouseState(mouseState);
    };

    const touch = new Guacamole.Mouse.Touchpad(element); // or Guacamole.Touchscreen

    touch.onmousedown = touch.onmousemove = touch.onmouseup = function (state) {
      client.sendMouseState(state);
    };

    setGuacd({
      client,
      sink,
    });
  }

  useEffect(() => {
    let resize = debounce(() => {
      onWindowResize();
    });
    window.addEventListener('resize', resize);
    window.addEventListener('beforeunload', handleUnload);
    window.addEventListener('focus', handleWindowFocus);

    return () => {
      window.removeEventListener('resize', resize);
      window.removeEventListener('beforeunload', handleUnload);
      window.removeEventListener('focus', handleWindowFocus);
    };
  }, [guacd])

  const onWindowResize = () => {
    if (guacd.client && !fixedSize) {
      const display = guacd.client.getDisplay();
      let width = window.innerWidth;
      let height = window.innerHeight;
      setBox({width, height});
      let scale = Math.min(
        height / display.getHeight(),
        width / display.getHeight()
      );
      display.scale(scale);
      guacd.client.sendSize(width, height);
    }
  }

  const handleUnload = (e) => {
    const message = "Want to leave the website?";
    (e || window.event).returnValue = message; //Gecko + IE
    return message;
  }

  const focus = () => {
    if (guacd.sink) {
      guacd.sink.focus();
    }
  }

  const handleWindowFocus = (e) => {
    if (navigator.clipboard) {
      try {
        navigator.clipboard.readText().then((text) => {
          sendClipboard({
            'data': text,
            'type': 'text/plain'
          });
        })
      } catch (e) {
        console.error('Copying to clipboard failed', e);
      }
    }
  };

  const handleClipboardReceived = (stream, mimetype) => {
    if (session['copy'] === '0') {
      return
    }

    if (/^text\//.exec(mimetype)) {
      let reader = new Guacamole.StringReader(stream);
      let data = '';
      reader.ontext = function textReceived(text) {
        data += text;
      };
      reader.onend = async () => {
        setClipboardText(data);
        if (navigator.clipboard) {
          await navigator.clipboard.writeText(data);
        }
      };
    } else {
      let reader = new Guacamole.BlobReader(stream, mimetype);
      reader.onend = () => {
        setClipboardText(reader.getBlob());
      }
    }
  };

  const sendClipboard = (data) => {
    if (!guacd.client) {
      return;
    }
    if (session['paste'] === '0') {
      message.warn('Can not paste');
      return
    }
    const stream = guacd.client.createClipboardStream(data.type);
    if (typeof data.data === 'string') {
      let writer = new Guacamole.StringWriter(stream);
      writer.sendText(data.data);
      writer.sendEnd();
    } else {
      let writer = new Guacamole.BlobWriter(stream);
      writer.oncomplete = function clipboardSent() {
        writer.sendEnd();
      };
      writer.sendBlob(data.data);
    }

    if (data.data && data.data.length > 0) {
    }
  }

  const onClientStateChange = (state, sessionId) => {
    const key = 'message';
    switch (state) {
      case STATE_IDLE:
        message.destroy(key);
        message.loading({content: 'Initializing...', duration: 0, key: key});
        break;
      case STATE_CONNECTING:
        message.destroy(key);
        message.loading({content: 'Connecting...', duration: 0, key: key});
        break;
      case STATE_WAITING:
        message.destroy(key);
        message.loading({content: 'Waiting for server response...', duration: 0, key: key});
        break;
      case STATE_CONNECTED:
        Modal.destroyAll();
        message.destroy(key);
        message.success({content: 'Connection successful', duration: 3, key: key});
        // Send a request to the backend to update the session's status
        // SessionBackend.connect(sessionId);
        break;
      case STATE_DISCONNECTING:
        // Handle disconnecting state if needed
        break;
      case STATE_DISCONNECTED:
        message.error({content: 'Connection closed', duration: 3, key: key});
        break;
      default:
        break;
    }
  };

  const sendCombinationKey = (keys) => {
    if (!guacd.client) {
      return;
    }
    for (let i = 0; i < keys.length; i++) {
      guacd.client.sendKeyEvent(1, keys[i]);
    }
    for (let j = 0; j < keys.length; j++) {
      guacd.client.sendKeyEvent(0, keys[j]);
    }
    message.success('Combination Keys successfully sent');
  }

  const showMessage = (msg) => {
    message.destroy();
    Modal.confirm({
      title: `Failed to connect to: ${assetName}`,
      icon: <CloseCircleOutlined />,
      content: msg,
      centered: true,
      okText: 'Reconnect',
      cancelText: 'Close this page',
      cancelButtonProps: {'danger': true},
      onOk() {
        window.location.reload();
      },
      onCancel() {
        window.close();
      },
    });
  }

  const onError = (status) => {
    switch (status.code) {
      case 256:
        showMessage('Unsupported access');
        break;
      case 512:
        showMessage('Remote service exception, please check if the target device can be accessed normally.');
        break;
      case 513:
        showMessage('Server busy');
        break;
      case 514:
        showMessage('Server connection timed out');
        break;
      case 515:
        showMessage('Remote service exception');
        break;
      case 516:
        showMessage('Resource not found');
        break;
      case 517:
        showMessage('Resource conflict');
        break;
      case 518:
        showMessage('Resource closed');
        break;
      case 519:
        showMessage('Remote service not found');
        break;
      case 520:
        showMessage('Remote service unavailable');
        break;
      case 521:
        showMessage('Session conflict');
        break;
      case 522:
        showMessage('Session connection timed out');
        break;
      case 523:
        showMessage('Session closed');
        break;
      case 768:
        showMessage('Network unreachable');
        break;
      case 769:
        showMessage('Server password authentication failed');
        break;
      case 771:
        showMessage('Client is forbidden');
        break;
      case 776:
        showMessage('Client connection timed out');
        break;
      case 781:
        showMessage('Client exception');
        break;
      case 783:
        showMessage('Incorrect request type');
        break;
      case 800:
        showMessage('Session does not exist');
        break;
      case 801:
        showMessage('Failed to create tunnel, please check if Guacd service is functioning correctly.');
        break;
      case 802:
        showMessage('Admin forcefully closed this session');
        break;
      default:
        if (status.message) {
          showMessage(Base64.decode(status.message));
        } else {
          showMessage('Unknown error.');
        }
    }
  };

  const fullScreen = () => {
    if (fullScreened) {
      exitFull();
      setFullScreened(false);
    } else {
      requestFullScreen(document.documentElement);
      setFullScreened(true);
    }
    focus();
  }

  const hotKeyMenu = (
    <Menu>
      <Menu.Item key={'ctrl+alt+delete'} onClick={() => sendCombinationKey(['65507', '65513', '65535'])}>Ctrl+Alt+Delete</Menu.Item>
      <Menu.Item key={'ctrl+alt+backspace'} onClick={() => sendCombinationKey(['65507', '65513', '65288'])}>Ctrl+Alt+Backspace</Menu.Item>
      <Menu.Item key={'windows+d'} onClick={() => sendCombinationKey(['65515', '100'])}>Windows+D</Menu.Item>
      <Menu.Item key={'windows+e'} onClick={() => sendCombinationKey(['65515', '101'])}>Windows+E</Menu.Item>
      <Menu.Item key={'windows+r'} onClick={() => sendCombinationKey(['65515', '114'])}>Windows+R</Menu.Item>
      <Menu.Item key={'windows+x'} onClick={() => sendCombinationKey(['65515', '120'])}>Windows+X</Menu.Item>
      <Menu.Item key={'windows'} onClick={() => sendCombinationKey(['65515'])}>Windows</Menu.Item>
    </Menu>
  );

  return (
    <div>
      <div className="container" style={{
        width: box.width,
        height: box.height,
        margin: '0 auto',
        backgroundColor: '#1b1b1b'
      }}>
        <div id="display"/>
      </div>
      <Draggable>
        <Affix style={{position: 'absolute', top: 50, right: 50}}>
          <Button icon={<ExpandOutlined/>} onClick={() => {
            fullScreen();
          }}/>
        </Affix>
      </Draggable>
      {
        session['copy'] === '1' || session['paste'] === '1' ?
          <Draggable>
            <Affix style={{position: 'absolute', top: 50, right: 100}}>
              <Button icon={<CopyOutlined/>}
                  onClick={() => {
                    setClipboardVisible(true);
                  }}/>
            </Affix>
          </Draggable> : undefined
      }
      {
        protocol === 'vnc' &&
        <Draggable>
          <Affix style={{position: 'absolute', top: 100, right: 100}}>
            <Dropdown overlay={hotKeyMenu} trigger={['click']} placement="bottomLeft">
              <Button icon={<WindowsOutlined/>}/>
            </Dropdown>
          </Affix>
        </Draggable>
      }
      {
        protocol === 'rdp' &&
        <Draggable>
          <Affix style={{position: 'absolute', top: 100, right: 100}}>
            <Dropdown overlay={hotKeyMenu} trigger={['click']} placement="bottomLeft">
              <Button icon={<WindowsOutlined/>}/>
            </Dropdown>
          </Affix>
        </Draggable>
      }
      <GuacdClipboard visible={clipboardVisible}
        clipboardText={clipboardText}
        handleOk={(text) => {
          sendClipboard({
            'data': text,
            'type': 'text/plain'
          });
          setClipboardText(text);
          setClipboardVisible(false);
          focus();
        }}
        handleCancel={() => {
          setClipboardVisible(false);
          focus();
        }}
      />
    </div>
  );
};

export default GuacdPage;
