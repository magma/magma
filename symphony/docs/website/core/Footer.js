/**
 * Copyright (c) 2017-present, Facebook, Inc.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

const React = require('react');

class Footer extends React.Component {
  docUrl(doc) {
    const baseUrl = this.props.config.baseUrl;
    const docsUrl = this.props.config.docsUrl;
    const docsPart = `${docsUrl ? `${docsUrl}/` : ''}`;
    return `${baseUrl}${docsPart}${doc}`;
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
                height="58"
              />
            )}
          </a>
          <div>
            <h5>Inventory Docs</h5>
            <a href={this.docUrl('csv-upload.html')}>
              CSV uploads
            </a>
            <a href={this.docUrl('py-inventory.html')}>
              python API
            </a>
            <a href={this.docUrl('export.html')}>
              Export
            </a>
          </div>
          <div>
            <h5>NMS Docs</h5>
            <a href={this.docUrl('nms-overview.html')}>
              Overview
            </a>
          </div>
          <div>
            <h5>Go to</h5>
              <a href={"/"}>
                  Inventory Management
              </a>
              <a href={"/workorders"}>
                  Workforce  Management
              </a>
              <a href={"/nms"}>
                  Network Management
              </a>
          </div>

        </section>

        <a
          href="https://opensource.facebook.com/"
          target="_blank"
          rel="noreferrer noopener"
          className="fbOpenSource">
          <img
            src={`${this.props.config.baseUrl}img/oss_logo.png`}
            alt="Facebook Open Source"
            width="170"
            height="45"
          />
        </a>
        <section className="copyright">{this.props.config.copyright}</section>
      </footer>
    );
  }
}

module.exports = Footer;
