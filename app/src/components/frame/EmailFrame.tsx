import React from 'react';
import {
  UiInputProps,
  UiInputDescriptionProps,
  UiInputCautionProps,
  UiInputLabelProps,
  UiInputLabel,
  UiInputDescription,
  UiInputCaution,
  UiInput,
} from 'components/atoms/input';
import { NonNullablePick } from 'components';

// email components
export type EmailFrameProps = {
  input: NonNullablePick<UiInputProps, 'value' | 'onChange' | 'readOnly'>;
  description: Pick<UiInputDescriptionProps, 'value'>;
  coution: Pick<UiInputCautionProps, 'value'>;
};
export const EmailFrame = (props: EmailFrameProps) => {
  const labelProps: UiInputLabelProps = {
    value: 'Email',
  };
  const inputProps: UiInputProps = {
    type: 'email',
    ...props.input,
  };
  const descriptionProps: UiInputDescriptionProps = {
    ...props.description,
  };
  const coutionProps: UiInputDescriptionProps = {
    ...props.coution,
  };
  return (
    <>
      <UiInputLabel {...labelProps} />
      <UiInputDescription {...descriptionProps} />
      <UiInputCaution {...coutionProps} />
      <UiInput {...inputProps} />
    </>
  );
};
