import React from 'react';
import { VerificationCodeForm, VerificationCodeFormProps } from 'components/frame/VerificationCodeForm';
import { UserBasicSigninForm, UserBasicSigninFormProps } from './UserBasicSigninForm';
import { UserNameForm, UserNameFormProps } from './UserNameForm';
import {
  SwitchAnimContent,
  SwitchAnimContentProps,
  SwitchAnimContext,
  SwitchAnimContextProps,
} from 'components/anims/SwitchAnim';

type AnimType = 'basic' | 'code' | 'name' | '';

export type SigninFormsProps = {
  switchContext: SwitchAnimContextProps<AnimType>;
  basic: UserBasicSigninFormProps;
  verificationCode: VerificationCodeFormProps;
  userName: UserNameFormProps;
};

export const SigninForms = ({ basic, verificationCode, userName, switchContext }: SigninFormsProps) => {
  const userBasicSigninFormProps: UserBasicSigninFormProps = {
    ...basic,
  };

  const verificationCodeFormProps: VerificationCodeFormProps = {
    ...verificationCode,
  };

  const userNameFormProps: UserNameFormProps = {
    ...userName,
  };

  const switchAnimContextProps: SwitchAnimContextProps<AnimType> = {
    ...switchContext,
  };
  const basicContentProps: SwitchAnimContentProps<AnimType> = {
    content: 'basic',
  };
  const codeContentProps: SwitchAnimContentProps<AnimType> = {
    content: 'code',
  };
  const nameContentProps: SwitchAnimContentProps<AnimType> = {
    content: 'name',
  };

  return (
    <SwitchAnimContext {...switchAnimContextProps}>
      <SwitchAnimContent {...basicContentProps}>
        <UserBasicSigninForm {...userBasicSigninFormProps} />
      </SwitchAnimContent>
      <SwitchAnimContent {...codeContentProps}>
        <VerificationCodeForm {...verificationCodeFormProps} />
      </SwitchAnimContent>
      <SwitchAnimContent {...nameContentProps}>
        <UserNameForm {...userNameFormProps} />
      </SwitchAnimContent>
    </SwitchAnimContext>
  );
};
