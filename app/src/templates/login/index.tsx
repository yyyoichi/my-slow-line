import React from 'react';
import { UiHead, UiMain, UiYCenter } from 'components/frame';
import { LoginForms } from 'components/login';
import { useLoginFormsProps } from './loginFormsProps';

export default function Login() {
  const props = useLoginFormsProps();
  return (
    <UiMain>
      <UiHead />
      <UiYCenter>
        <LoginForms {...props} />
      </UiYCenter>
    </UiMain>
  );
}
