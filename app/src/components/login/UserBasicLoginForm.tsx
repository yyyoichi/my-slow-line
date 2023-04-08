import { EmailFrameProps, EmailFrame } from 'components/frame/EmailFrame';
import { LoadButton, LoadButtonProps } from 'components/frame/LoadButton';
import { PasswordFrameProps, PasswordFrame } from 'components/frame/PasswordFrame';
import FormFrame from 'components/frame/FormFrame';
import React from 'react';

export type UserBasicLoginFormProps = {
  password: PasswordFrameProps;
  email: EmailFrameProps;
  sendButton: SendButtonProps;
};

export const UserBasicLoginForm = (props: UserBasicLoginFormProps) => {
  return (
    <>
      <FormFrame.Container>
        <FormFrame.Content>
          <EmailFrame {...props.email} />
        </FormFrame.Content>
        <FormFrame.Content>
          <PasswordFrame {...props.password} />
        </FormFrame.Content>
        <div className='mx-auto w-1/3  py-4'>
          <SendButton {...props.sendButton} />
        </div>
      </FormFrame.Container>
    </>
  );
};

// send button

type SendButtonProps = Pick<LoadButtonProps, 'active' | 'onClick'>;

const SendButton = (props: SendButtonProps) => {
  const loadButtonProps: LoadButtonProps = {
    color: 'yellow',
    ...props,
  };
  return <LoadButton {...loadButtonProps}>{'login'}</LoadButton>;
};
