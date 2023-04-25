import { EmailFrameProps, EmailFrame } from 'components/frame/EmailFrame';
import { LoadButton, LoadButtonProps } from 'components/frame/LoadButton';
import FormFrame from 'components/frame/FormFrame';
import { PasswordFrameProps, PasswordFrame, PasswordFieldProps, PasswordField } from 'components/frame/PasswordFrame';
import React from 'react';
import { NonNullablePick } from 'components';

export type UserBasicSigninFormProps = {
  password: PasswordFrameProps;
  confirmPassword: PasswordFieldProps;
  email: EmailFrameProps;
  sendButton: SendButtonProps;
};

export const UserBasicSigninForm = (props: UserBasicSigninFormProps) => {
  return (
    <>
      <FormFrame.Container>
        <h2 className='text-center text-3xl'>
          WELCOME TO
          <div className='py-2 font-extrabold text-my-yellow shadow-sm'>Ctrl+</div>
        </h2>
        <FormFrame.Content>
          <EmailFrame {...props.email} />
        </FormFrame.Content>
        <FormFrame.Content>
          <PasswordFrame {...props.password} />
          <div className='my-2'>
            <PasswordField {...props.confirmPassword} />
          </div>
        </FormFrame.Content>
        <div className='mx-auto w-1/3  py-4'>
          <SendButton {...props.sendButton} />
        </div>
      </FormFrame.Container>
    </>
  );
};

// send button

type SendButtonProps = NonNullablePick<LoadButtonProps, 'active' | 'onClick'>;

const SendButton = (props: SendButtonProps) => {
  const loadButtonProps: LoadButtonProps = {
    color: 'yellow',
    ...props,
  };
  return <LoadButton {...loadButtonProps}>{'Signin'}</LoadButton>;
};
