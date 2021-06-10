/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Ref: https://v1.docusaurus.io/docs/en/site-config

const url = process.env.DOCUSAURUS_URL || 'https://magmacore.org'
const baseUrl = process.env.DOCUSAURUS_BASE_URL || '/'

const mermaid = require('remark-mermaid')

const siteConfig = {
  title: 'Magma Documentation',
  disableTitleTagline: true,
  tagline: 'Bring more people online by enabling operators with open, flexible, and extensible network solutions',

  projectName: 'magma',
  organizationName: 'magma',
  url: url, // full URL
  baseUrl: baseUrl, // base URL

  headerLinks: [
    {href: 'https://magmacore.org', label: 'Home'},
    {label: ' | '},
    {href: baseUrl, label: 'Docs'},
    {label: ' | '},
    {href: 'https://github.com/magma', label: 'Code'},
    {label: ' | '},
    {href: 'https://magmacore.org/community', label: 'Community'},
  ],

  // Path to images for header/footer
  headerIcon: 'img/magma-logo.svg',
  footerIcon: 'img/magma_icon.png',
  favicon: 'img/icon.png',

  // Colors for website
  colors: {
    primaryColor: '#5602a4',
    secondaryColor: '#5602a4',
  },

  // This copyright info is used in /core/Footer.js and blog RSS/Atom feeds.
  copyright: `Copyright \u{00A9} ${new Date().getFullYear()} The Magma Authors`,

  // Highlight.js theme to use for syntax highlighting in code blocks.
  highlight: {
    theme: 'default',
  },

  // Add custom scripts here that would be placed in <script> tags.
  scripts: ['https://buttons.github.io/buttons.js',
    'https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js',
    '/init.js'],

  // On page navigation for the current documentation page.
  onPageNav: 'separate',
  // No .html extensions for paths.
  cleanUrl: true,

  // Open Graph and Twitter card images.
  ogImage: 'img/docusaurus.png',
  twitterImage: 'img/docusaurus.png',

  // Whether sidebar nav is collapsible
  docsSideNavCollapsible: true,

  // Enable Algolia DocSearch Functionality within Docusaurus
  algolia: {
    apiKey: 'f95caeb7bc059b294eec88e340e5445b',
    indexName: 'magma',
  },

  // Enable mermaid
  markdownPlugins: [ (md) => {
        md.renderer.rules.fence_custom.mermaid = (tokens, idx, options, env, instance) => {
            return `<div class="mermaid">${tokens[idx].content}</div>`;
        };
    }],



};

module.exports = siteConfig;
