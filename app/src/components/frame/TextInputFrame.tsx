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

// email components
export type TextInputFrameProps = {
  label: Pick<UiInputLabelProps, 'value'>;
  input: Pick<UiInputProps, 'value' | 'onChange' | 'readOnly'>;
  description: Pick<UiInputDescriptionProps, 'value'>;
  coution: Pick<UiInputCautionProps, 'value'>;
};
export const TextInputFrame = (props: TextInputFrameProps) => {
  const labelProps: UiInputLabelProps = {
    ...props.label,
  };
  const inputProps: UiInputProps = {
    type: 'text',
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
