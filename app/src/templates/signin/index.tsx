import React from 'react';
import { UiHead, UiMain, UiYCenter } from 'components/frame';
import SigninForms from './signinForms';

export default function Signin() {
  return (
    <UiMain>
      <UiHead />
      <UiYCenter>
        <SigninForms />
      </UiYCenter>
    </UiMain>
  );
}
