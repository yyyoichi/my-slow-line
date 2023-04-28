import React from 'react';
import { VerificationCodeForm, VerificationCodeFormProps } from 'components/frame/VerificationCodeForm';
import { UserBasicLoginForm, UserBasicLoginFormProps } from './UserBasicLoginForm';
import {
  SwitchAnimContent,
  SwitchAnimContentProps,
  SwitchAnimContext,
  SwitchAnimContextProps,
} from 'components/anims/SwitchAnim';

type AnimType = 'basic' | 'code' | '';

export type LoginFormsProps = {
  switchContext: SwitchAnimContextProps<AnimType>;
  basic: UserBasicLoginFormProps;
  verificationCode: VerificationCodeFormProps;
};

export const LoginForms = ({ basic, verificationCode, switchContext }: LoginFormsProps) => {
  const userBasicLoginFormProps: UserBasicLoginFormProps = {
    ...basic,
  };

  const verificationCodeFormProps: VerificationCodeFormProps = {
    ...verificationCode,
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

  return (
    <SwitchAnimContext {...switchAnimContextProps}>
      <SwitchAnimContent {...basicContentProps}>
        <UserBasicLoginForm {...userBasicLoginFormProps} />
      </SwitchAnimContent>
      <SwitchAnimContent {...codeContentProps}>
        <VerificationCodeForm {...verificationCodeFormProps} />
      </SwitchAnimContent>
    </SwitchAnimContext>
  );
};
