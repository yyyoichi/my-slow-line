import 'styles/globals.css';
import React from 'react';
import Head from 'next/head';
import type { AppProps } from 'next/app';

export default function App({ Component, pageProps }: AppProps) {
  return (
    <>
      <Head>
        <link rel='icon' sizes='16x16' href='/favicon.ico' />
      </Head>
      <Component {...pageProps} />
    </>
  );
}
