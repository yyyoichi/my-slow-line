import React from 'react';
import Head from 'next/head';
import type { AppProps } from 'next/app';
import { MyAccountContext, useMyAccountState } from 'hooks';
import { FullScreentitle } from 'components/frame/FullScreen';
import { Backscreen } from 'components/atoms/backscreen';

export default function App({ Component, pageProps }: AppProps) {
  const context = useMyAccountState();
  return (
    <>
      <Head>
        <link rel='icon' sizes='16x16' href='/favicon.ico' />
      </Head>
      <Backscreen>
        <FullScreentitle active={!context.initialized} />
        {context.initialized && (
          <MyAccountContext.Provider value={context}>
            <Component {...pageProps} />
          </MyAccountContext.Provider>
        )}
      </Backscreen>
    </>
  );
}
