import { UiMain, UiYCenter } from 'components/frame';
import React from 'react';
import { Head } from 'components/head';
import { useHeadProps } from './headProps';
import { useHomePageProps } from './homePageProps';
import { HomePage } from 'components/home';
export default function Home() {
  const headProps = useHeadProps();
  const homeProps = useHomePageProps();
  return (
    <UiMain>
      <Head {...headProps} />
      <UiYCenter>
        <HomePage {...homeProps} />
      </UiYCenter>
    </UiMain>
  );
}
