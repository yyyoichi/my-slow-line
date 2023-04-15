import React from 'react';
import { UiHead, UiMain, UiYCenter } from 'components/frame';
import LoginForms from './loginForms';

export default function Login() {
  return (
    <UiMain>
      <UiHead />
      <UiYCenter>
        <LoginForms />
      </UiYCenter>
    </UiMain>
  );
}
