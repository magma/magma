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

const React = require('react');

class Footer extends React.Component {
  docUrl(doc, language) {
    const baseUrl = this.props.config.baseUrl;
    const docsUrl = this.props.config.docsUrl;
    const docsPart = `${docsUrl ? `${docsUrl}/` : ''}`;
    const langPart = `${language ? `${language}/` : ''}`;
    return `${baseUrl}${docsPart}${langPart}${doc}`;
  }

  pageUrl(doc, language) {
    const baseUrl = this.props.config.baseUrl;
    return baseUrl + (language ? `${language}/` : '') + doc;
  }

  render() {
    return (
      <footer className="nav-footer" id="footer">
        <section className="sitemap">
          <a href={this.props.config.baseUrl} className="nav-home">
            {this.props.config.footerIcon && (
              <img
                src={this.props.config.baseUrl + this.props.config.footerIcon}
                alt={this.props.config.title}
                width="66"
              />
            )}
          </a>
          <div>
            <h5>Docs</h5>
            <a href="https://github.com/facebookincubator/magma/blob/master/docs/Magma_Product_Overview.pdf">
              Magma Product Overview
            </a>
            <a href="https://github.com/facebookincubator/magma/blob/master/docs/Magma_Specs_V1.1.pdf">
              Magma Spec
            </a>
          </div>
          <div>
            <h5>Community</h5>
            <a href="https://discord.gg/4YxZbft">
              Discord
            </a>

            <a
              href="https://fb.me/magmadevsummit"
              target="_blank"
              rel="noreferrer noopener">
              Magma Dev Summit
            </a>
          </div>
          <div>
            <h5>More</h5>
            <a href="https://code.fb.com/open-source/magma/">Blog</a>
            <a href="https://github.com/facebookincubator/magma">GitHub</a>
          </div>
        </section>

        <section className="copyright">{this.props.config.copyright}</section>
      </footer>
    );
  }
}

module.exports = Footer;
