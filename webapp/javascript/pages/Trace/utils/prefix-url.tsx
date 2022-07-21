// Copyright (c) 2017 Uber Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// import sitePrefix from '../site-prefix';

const baseNode = document.querySelector('base');
if (!baseNode && process.env.NODE_ENV !== 'test') {
  throw new Error('<base> element not found');
}

const sitePrefix = baseNode ? baseNode.href : `${global.location.origin}/`;

// Configure the webpack publicPath to match the <base>:
// https://webpack.js.org/guides/public-path/#on-the-fly
// eslint-disable-next-line camelcase
window.__webpack_public_path__ = sitePrefix;

// export default sitePrefix;

const origin =
  process.env.NODE_ENV === 'test'
    ? global.location.origin
    : window.location.origin;

export function getPathPrefix(orig?: string, sitePref?: string) {
  const o = orig == null ? '' : orig.replace(/[-/\\^$*+?.()|[\]{}]/g, '\\$&');
  const s = sitePref == null ? '' : sitePref;
  const rx = new RegExp(`^${o}|/$`, 'ig');
  return s.replace(rx, '');
}

const pathPrefix = getPathPrefix(origin, sitePrefix);

export default function prefixUrl(value?: string) {
  const s = value == null ? '' : String(value);
  return `${pathPrefix}${s}`;
}
