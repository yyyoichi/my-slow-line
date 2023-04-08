import React from 'react';
import { UiHead, UiMain, UiYCenter } from 'components/frame';
import { UserBasicSigninForm, UserBasicSigninFormProps } from 'components/signin/UserBasicSigninForm';
import { UserNameForm, UserNameFormProps } from 'components/signin/UserNameForm';
import { VerificationCodeForm, VerificationCodeFormProps } from 'components/frame/VerificationCodeForm';

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

  const verificationCodeFormProps: VerificationCodeFormProps = {
    name: {
      input: {
        value: '000000',
        onChange: () => console.log('click'),
        readOnly: false,
      },
      description: {
        value: '6 code',
      },
      coution: {
        value: '',
      },
    },
    sendButton: {
      active: false,
      onClick: () => console.log('click'),
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
        <VerificationCodeForm {...verificationCodeFormProps} />
        <UserNameForm {...userNameFormProps} />
      </UiYCenter>
    </UiMain>
  );
}
