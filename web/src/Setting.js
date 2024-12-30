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

import React from "react";
import {Tooltip, message} from "antd";
import {isMobile as isMobileDevice} from "react-device-detect";
import i18next from "i18next";
import Sdk from "casdoor-js-sdk";
import {QuestionCircleTwoTone} from "@ant-design/icons";
import {v4 as uuidv4} from "uuid";

export let ServerUrl = "";
export let CasdoorSdk;

export function initServerUrl() {
  const hostname = window.location.hostname;
  if (hostname === "localhost") {
    ServerUrl = `http://${hostname}:19000`;
  }
}

export function isLocalhost() {
  const hostname = window.location.hostname;
  return hostname === "localhost";
}

export function initCasdoorSdk(config) {
  CasdoorSdk = new Sdk(config);
}

function getUrlWithLanguage(url) {
  if (url.includes("?")) {
    return `${url}&language=${getLanguage()}`;
  } else {
    return `${url}?language=${getLanguage()}`;
  }
}

export function getSignupUrl() {
  return getUrlWithLanguage(CasdoorSdk.getSignupUrl());
}

export function getSigninUrl() {
  return getUrlWithLanguage(CasdoorSdk.getSigninUrl());
}

export function getUserProfileUrl(userName, account) {
  return getUrlWithLanguage(CasdoorSdk.getUserProfileUrl(userName, account));
}

export function getMyProfileUrl(account) {
  return getUrlWithLanguage(CasdoorSdk.getMyProfileUrl(account));
}

export function signin() {
  return CasdoorSdk.signin(ServerUrl);
}

export function parseJson(s) {
  if (s === "") {
    return null;
  } else {
    return JSON.parse(s);
  }
}

export function myParseInt(i) {
  const res = parseInt(i);
  return isNaN(res) ? 0 : res;
}

export function openLink(link) {
  // this.props.history.push(link);
  const w = window.open("about:blank");
  w.location.href = link;
}

export function goToLink(link) {
  window.location.href = link;
}

export function goToLinkSoft(ths, link) {
  ths.props.history.push(link);
}

export function showMessage(type, text) {
  if (type === "") {
    return;
  } else if (type === "success") {
    message.success(text);
  } else if (type === "error") {
    message.error(text);
  }
}

export function isAdminUser(account) {
  return account?.isAdmin;
}

export function deepCopy(obj) {
  return Object.assign({}, obj);
}

export function insertRow(array, row, i) {
  return [...array.slice(0, i), row, ...array.slice(i)];
}

export function addRow(array, row) {
  return [...array, row];
}

export function prependRow(array, row) {
  return [row, ...array];
}

export function deleteRow(array, i) {
  // return array = array.slice(0, i).concat(array.slice(i + 1));
  return [...array.slice(0, i), ...array.slice(i + 1)];
}

export function swapRow(array, i, j) {
  return [...array.slice(0, i), array[j], ...array.slice(i + 1, j), array[i], ...array.slice(j + 1)];
}

export function isMobile() {
  // return getIsMobileView();
  return isMobileDevice;
}

export function getFormattedDate(date) {
  if (date === undefined || date === null) {
    return null;
  }

  date = date.replace("T", " ");
  date = date.replace("+08:00", " ");
  return date;
}

export function getFormattedDateShort(date) {
  return date.slice(0, 10);
}

export function getShortName(s) {
  return s.split("/").slice(-1)[0];
}

export function getShortText(s, maxLength = 35) {
  if (s.length > maxLength) {
    return `${s.slice(0, maxLength)}...`;
  } else {
    return s;
  }
}

export function getRandomName() {
  return Math.random().toString(36).slice(-6);
}

function getRandomInt(s) {
  let hash = 0;
  if (s.length !== 0) {
    for (let i = 0; i < s.length; i++) {
      const char = s.charCodeAt(i);
      hash = ((hash << 5) - hash) + char;
      hash = hash & hash;
    }
  }

  return hash;
}

export function getAvatarColor(s) {
  const colorList = ["#f56a00", "#7265e6", "#ffbf00", "#00a2ae"];
  let random = getRandomInt(s);
  if (random < 0) {
    random = -random;
  }
  return colorList[random % 4];
}

export function getLanguage() {
  return i18next.language;
}

export function setLanguage(language) {
  localStorage.setItem("language", language);
  changeMomentLanguage(language);
  i18next.changeLanguage(language);
}

export function changeLanguage(language) {
  localStorage.setItem("language", language);
  changeMomentLanguage(language);
  i18next.changeLanguage(language);
  window.location.reload(true);
}

export function changeMomentLanguage(lng) {
  return;
  // if (lng === "zh") {
  //   moment.locale("zh", {
  //     relativeTime: {
  //       future: "%s内",
  //       past: "%s前",
  //       s: "几秒",
  //       ss: "%d秒",
  //       m: "1分钟",
  //       mm: "%d分钟",
  //       h: "1小时",
  //       hh: "%d小时",
  //       d: "1天",
  //       dd: "%d天",
  //       M: "1个月",
  //       MM: "%d个月",
  //       y: "1年",
  //       yy: "%d年",
  //     },
  //   });
  // }
}

export function getProviderTypeOptions(category) {
  if (category === "Public Cloud") {
    return ([
      {value: "Amazon Web Services", label: "Amazon Web Services"},
      {value: "Azure", label: "Azure"},
      {value: "Google Cloud", label: "Google Cloud"},
      {value: "Aliyun", label: "Aliyun"},
    ]);
  } else if (category === "Private Cloud") {
    return ([
      {value: "KVM", label: "KVM"},
      {value: "Xen", label: "Xen"},
      {value: "VMware", label: "VMware"},
      {value: "PVE", label: "PVE"},
    ]);
  } else if (category === "Blockchain") {
    return ([
      {value: "Hyperledger Fabric", label: "Hyperledger Fabric"},
      {value: "ChainMaker", label: "ChainMaker"},
    ]);
  } else {
    return [];
  }
}

export function setOrganization(organization) {
  localStorage.setItem("organization", organization);
  window.dispatchEvent(new Event("storageOrganizationChanged"));
}

export function isDefaultOrganizationSelected(account) {
  if (isAdminUser(account)) {
    return getOrganization() === "All";
  }
  return false;
}

export function getOrganization() {
  const organization = localStorage.getItem("organization");
  return organization !== null ? organization : "All";
}

export function getRequestOrganization(account) {
  if (isAdminUser(account)) {
    return getOrganization() === "All" ? account.owner : getOrganization();
  }
  return account.owner;
}

export function isLocalAdminUser(account) {
  if (account === undefined || account === null) {
    return false;
  }
  return account.isAdmin === true || isAdminUser(account);
}

export function getAcceptLanguage() {
  if (i18next.language === null || i18next.language === "") {
    return "en;q=0.9,en;q=0.8";
  }
  return i18next.language + ";q=0.9,en;q=0.8";
}

export const StaticBaseUrl = "https://cdn.casbin.org";

export const Countries = [{label: "English", key: "en", country: "US", alt: "English"},
  {label: "Español", key: "es", country: "ES", alt: "Español"},
  {label: "Français", key: "fr", country: "FR", alt: "Français"},
  {label: "Deutsch", key: "de", country: "DE", alt: "Deutsch"},
  {label: "中文", key: "zh", country: "CN", alt: "中文"},
  {label: "Indonesia", key: "id", country: "ID", alt: "Indonesia"},
  {label: "日本語", key: "ja", country: "JP", alt: "日本語"},
  {label: "한국어", key: "ko", country: "KR", alt: "한국어"},
  {label: "Русский", key: "ru", country: "RU", alt: "Русский"},
  {label: "TiếngViệt", key: "vi", country: "VN", alt: "TiếngViệt"},
  {label: "Português", key: "pt", country: "BR", alt: "Português"},
  {label: "Itariano", key: "it", country: "IT", alt: "Itariano"},
  {label: "Marley", key: "ms", country: "MY", alt: "Marley"},
  {label: "Tkiš", key: "tr", country: "TR", alt: "Tkiš"},
  {label: "لغة عربية", key: "ar", country: "DZ", alt: "لغة عربية"},
  {label: "עִבְרִית", key: "he", country: "IL", alt: "עִבְרִית"},
  {label: "Filipino", key: "fi", country: "PH", alt: "Filipino"},
];

export function getLabel(text, tooltip) {
  return (
    <React.Fragment>
      <span style={{marginRight: 4}}>{text}</span>
      <Tooltip placement="top" title={tooltip}>
        <QuestionCircleTwoTone twoToneColor="rgb(45,120,213)" />
      </Tooltip>
    </React.Fragment>
  );
}

export function getItem(label, key, icon, children, type) {
  return {
    key,
    icon,
    children,
    label,
    type,
  };
}

export function getOption(label, value) {
  return {
    label,
    value,
  };
}

export function isResponseDenied(res) {
  if (res.msg === "Unauthorized operation" || res.msg === "未授权的操作") {
    return true;
  }
  return false;
}

export function GenerateId() {
  return uuidv4();
}

export function GetIdFromObject(obj) {
  if (obj === undefined || obj === null) {
    return "";
  }
  return `${obj.owner}/${obj.name}`;
}
