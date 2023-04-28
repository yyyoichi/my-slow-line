import { UiMain } from 'components/frame';
import React from 'react';
import { Head } from 'components/home';
import { useHeadProps } from './headProps';
export default function Home() {
  const headProps = useHeadProps();
  return (
    <UiMain>
      <Head {...headProps} />
    </UiMain>
  );
}
