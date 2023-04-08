import React from 'react';
import { UiHead, UiMain, UiYCenter } from 'components/frame';
import { UserBasicSigninForm, UserBasicSigninFormProps } from 'components/signin/UserBasicSigninForm';
import { UserNameForm, UserNameFormProps } from 'components/signin/UserNameForm';

export default function Signin() {
  const userBasicSigninFormProps: UserBasicSigninFormProps = {
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
      active: true,
      onClick: () => console.log('click'),
    },
    confirmPassword: {
      input: {
        value: 'confirm',
        onChange: () => console.log('r'),
        readOnly: true,
      },
    },
  };

  const userNameFormProps: UserNameFormProps = {
    name: {
      input: {
        value: 'name',
        onChange: () => console.log('click'),
        readOnly: false,
      },
      coution: {
        value: '',
      },
      description: {
        value: 'user account name',
      },
    },
    sendButton: {
      active: true,
      onClick: () => console.log('click'),
    },
  };

  return (
    <UiMain>
      <UiHead />
      <UiYCenter>
        <UserBasicSigninForm {...userBasicSigninFormProps} />
        <UserNameForm {...userNameFormProps} />
      </UiYCenter>
    </UiMain>
  );
}
