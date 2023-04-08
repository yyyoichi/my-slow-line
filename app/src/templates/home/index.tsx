import { UiHead, UiMain } from 'components/frame';
import Link from 'next/link';
import React from 'react';
export default function Home() {
  return (
    <UiMain>
      <UiHead></UiHead>
      <Link href={'/login'}>login</Link>
      <Link href={'/signin'}>signin</Link>
    </UiMain>
  );
}
