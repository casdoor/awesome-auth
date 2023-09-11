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

export const debounce = function (fn, delay = 500) {
  let timer = null;

  return function () {
    if (timer) {
      clearTimeout(timer)
    }
    timer = setTimeout(() => {
      fn.apply(this, arguments)
      timer = null
    }, delay)
  }
}

export const getToken = function () {
  return localStorage.getItem('X-Auth-Token');
}

export function requestFullScreen(element) {
  const requestMethod = element.requestFullScreen || //W3C
    element.webkitRequestFullScreen || //FireFox
    element.mozRequestFullScreen || //Chrome
    element.msRequestFullScreen; //IE11
  if (requestMethod) {
    requestMethod.call(element);
  } else if (typeof window.ActiveXObject !== "undefined") { //for Internet Explorer
    const wScript = new window.ActiveXObject("WScript.Shell");
    if (wScript !== null) {
      wScript.SendKeys("{F11}");
    }
  }
}

// exit full screen
export function exitFull() {
  const exitMethod = document.exitFullscreen || //W3C
    document.mozCancelFullScreen || //FireFox
    document.webkitExitFullscreen || //Chrome
    document.webkitExitFullscreen; //IE11
  if (exitMethod) {
    exitMethod.call(document);
  } else if (typeof window.ActiveXObject !== "undefined") { //for Internet Explorer
    const wScript = new window.ActiveXObject("WScript.Shell");
    if (wScript !== null) {
      wScript.SendKeys("{F11}");
    }
  }
}
