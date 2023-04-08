import React from 'react';
import { UiHead, UiMain, UiYCenter } from 'components/frame';
import { UserBasicLoginForm, UserBasicLoginFormProps } from 'components/login/UserBasicLoginForm';

export default function Login() {
  const userBasicLoginFormProps: UserBasicLoginFormProps = {
    password: {
      input: {
        value: 'text',
        onChange: (e) => {
          console.log(e.target.value);
        },
        readOnly: false,
      },
      description: {
        value: 'prese entry',
      },
      coution: {
        value: 'nil',
      },
    },
    email: {
      input: {
        value: 'text',
        onChange: (e) => {
          console.log(e.target.value);
        },
        readOnly: true,
      },
      description: {
        value: 'prese entry',
      },
      coution: {
        value: 'nil',
      },
    },
    sendButton: {
      active: false,
      onClick: () => console.log('click'),
    },
  };

  return (
    <UiMain>
      <UiHead />
      <UiYCenter>
        <UserBasicLoginForm {...userBasicLoginFormProps} />
      </UiYCenter>
    </UiMain>
  );
}
