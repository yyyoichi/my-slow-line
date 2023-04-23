import React from 'react';
import { UiHead, UiMain, UiYCenter } from 'components/frame';
import { useSigninFormsProps } from './signinFormsProps';
import { SigninForms } from 'components/signin';

export default function Signin() {
  const signinFormsProps = useSigninFormsProps();
  return (
    <UiMain>
      <UiHead />
      <UiYCenter>
        <SigninForms {...signinFormsProps} />
      </UiYCenter>
    </UiMain>
  );
}
