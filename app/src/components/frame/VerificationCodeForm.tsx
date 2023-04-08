import React from 'react';
import FormFrame from 'components/frame/FormFrame';
import { LoadButtonProps, LoadButton } from 'components/frame/LoadButton';
import { TextInputFrame, TextInputFrameProps } from 'components/frame/TextInputFrame';

export type VerificationCodeFormProps = {
  name: Pick<TextInputFrameProps, 'input' | 'description' | 'coution'>;
  sendButton: SendButtonProps;
};

export const VerificationCodeForm = (props: VerificationCodeFormProps) => {
  const textInputProps: TextInputFrameProps = {
    ...props.name,
    label: {
      value: 'Code',
    },
  };
  const buttonProps: SendButtonProps = {
    ...props.sendButton,
  };
  return (
    <FormFrame.Container>
      <h3 className='text-center text-2xl'>Check your email box</h3>
      <FormFrame.Content>
        <TextInputFrame {...textInputProps} />
        <div className='mx-auto my-3 w-1/3'>
          <SendButton {...buttonProps} />
        </div>
      </FormFrame.Content>
    </FormFrame.Container>
  );
};

// send button

type SendButtonProps = Pick<LoadButtonProps, 'active' | 'onClick'>;

const SendButton = (props: SendButtonProps) => {
  const loadButtonProps: LoadButtonProps = {
    color: 'yellow',
    ...props,
  };
  return <LoadButton {...loadButtonProps}>{'Verification'}</LoadButton>;
};
