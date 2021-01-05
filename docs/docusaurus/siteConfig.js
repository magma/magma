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

// See https://docusaurus.io/docs/site-config for all the possible
// site configuration options.

// List of projects/orgs using your project for the users page.
const users = [
  {
    caption: 'User1',
    // You will need to prepend the image path with your baseUrl
    // if it is not '/', like: '/test-site/img/docusaurus.svg'.
    image: '/img/magma-logo.svg',
    infoLink: 'https://magmacore.org',
    pinned: true,
  },
];

const url = process.env.DOCUSAURUS_URL || 'https://magmacore.org'
const baseUrl = process.env.DOCUSAURUS_BASE_URL || '/'

const siteConfig = {
  title: 'Magma Documentation', // Title for your website.
  disableTitleTagline: true,
  tagline: 'Bring more people online by enabling operators with open, flexible, and extensible network solutions',

  // Used for publishing and more
  projectName: 'magma',
  organizationName: 'magma',
  url: url, // Your website URL
  baseUrl: baseUrl, // Base URL for your project */
  // For github.io type URLs, you would set the url and baseUrl like:
  //   url: 'https://facebook.github.io',
  //   baseUrl: '/test-site/',
  // For top-level user or org sites, the organization is still the same.
  // e.g., for the https://JoelMarcey.github.io site, it would be set like...
  //   organizationName: 'Facebook'

  //docsUrl: 'docs',

  // For no header links in the top nav bar -> headerLinks: [],
  headerLinks: [
    // {doc: 'basics/introduction', label: 'Docs'},
    {href: 'https://magmacore.org', label: 'Home'},
    {label: ' | '},
    {href: '/', label: 'Docs'},
    {label: ' | '},
    {href: 'https://github.com/magma', label: 'Code'},
    {label: ' | '},
    {href: 'https://magmacore.org/community', label: 'Community'},
  ],

  // If you have users set above, you add it here:
  users,

  /* path to images for header/footer */
  headerIcon: 'img/magma-logo.svg',
  footerIcon: 'img/magma_icon.png',
  favicon: 'img/favicon.png',

  /* Colors for website */
  colors: {
    primaryColor: '#5602a4',
    secondaryColor: '#5602a4',
  },

  /* Custom fonts for website */
  /*
  fonts: {
    myFont: [
      "Times New Roman",
      "Serif"
    ],
    myOtherFont: [
      "-apple-system",
      "system-ui"
    ]
  },
  */

  // This copyright info is used in /core/Footer.js and blog RSS/Atom feeds.
  copyright: `Copyright \u{00A9} ${new Date().getFullYear()} The Magma Authors`,

  highlight: {
    // Highlight.js theme to use for syntax highlighting in code blocks.
    theme: 'default',
  },

  // Add custom scripts here that would be placed in <script> tags.
  scripts: ['https://buttons.github.io/buttons.js'],

  // On page navigation for the current documentation page.
  onPageNav: 'separate',
  // No .html extensions for paths.
  cleanUrl: true,

  // Open Graph and Twitter card images.
  ogImage: 'img/docusaurus.png',
  twitterImage: 'img/docusaurus.png',


  docsSideNavCollapsible: true,

  // Show documentation's last contributor's name.
  // enableUpdateBy: true,

  // Show documentation's last update time.
  // enableUpdateTime: true,

  // You may provide arbitrary config keys to be used as needed by your
  // template. For example, if you need your repo's URL...
  //   repoUrl: 'https://github.com/facebook/test-site',
};

module.exports = siteConfig;
