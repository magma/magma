import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import HomepageFeatures from '@site/src/components/HomepageFeatures';

import Heading from '@theme/Heading';
import styles from './index.module.css';

function HomepageHeader() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <div style={{justifyContent:'center', alignItems:'center', display:'flex', flexDirection:'column'}}>
        <h1 style={{color:'black', marginTop:'100px', marginBottom:'30px', fontSize:'2.5rem'}}>Welcome to Magma Documentation</h1>
        <img src='/img/icon.png' style={{marginBottom:'30px'}}></img>
        <div className={styles.buttons}>
          <Link
            className="button button--secondary button--lg"
            to="/docs/basics/version-1.8.0-introduction" >
            Go to Documentation
          </Link>
        </div>
    </div>
  );
}

export default function Home() {
  const {siteConfig} = useDocusaurusContext();
  return (
    <Layout
      title={`Hello from ${siteConfig.title}`}
      description="Description will go into a meta tag in <head />">
      <HomepageHeader />
    </Layout>
  );
}
