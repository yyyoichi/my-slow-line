import { EmailFrameProps, EmailFrame } from 'components/frame/EmailFrame';
import { LoadButton, LoadButtonProps } from 'components/frame/LoadButton';
import { PasswordFrameProps, PasswordFrame, PasswordFieldProps, PasswordField } from 'components/frame/PasswordFrame';
import React from 'react';

export type UserBasicSigninFormProps = {
  password: PasswordFrameProps;
  confirmPassword: PasswordFieldProps;
  email: EmailFrameProps;
  sendButton: SendButtonProps;
};

export const UserBasicSigninForm = (props: UserBasicSigninFormProps) => {
  const Wrap = ({ children }: { children: React.ReactNode }) => <div className='px-4 py-1'>{children}</div>;
  return (
    <>
      <div className='flex flex-col gap-3 rounded-md bg-my-light-black py-2 '>
        <h2 className='text-center text-3xl'>
          WELCOME TO
          <div className='py-2 font-extrabold text-my-yellow shadow-sm'>Ctrl+</div>
        </h2>
        <Wrap>
          <EmailFrame {...props.email} />
        </Wrap>
        <Wrap>
          <PasswordFrame {...props.password} />
          <div className='my-2'>
            <PasswordField {...props.confirmPassword} />
          </div>
        </Wrap>
        <div className='mx-auto w-1/3  py-4'>
          <SendButton {...props.sendButton} />
        </div>
      </div>
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
  return <LoadButton {...loadButtonProps}>{'Signin'}</LoadButton>;
};
