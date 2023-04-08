import React from 'react';
import { DivProps, InputProps } from './type';

export type UiInputProps = InputProps;
export const UiInput = ({ className = '', ...props }: UiInputProps) => (
  <input
    className={`block w-full rounded-sm bg-my-white px-2 text-lg text-my-black outline-none focus:outline-none ${className}`}
    {...props}
  />
);

export type UiInputLabelProps = {
  value: string;
} & DivProps;
export const UiInputLabel = ({ className = '', ...props }: UiInputLabelProps) => {
  return (
    <div className={`text-lg ${className}`} {...props}>
      {props.value}
    </div>
  );
};

export type UiInputDescriptionProps = {
  value: string;
} & DivProps;
export const UiInputDescription = ({ className = '', ...props }: UiInputDescriptionProps) => (
  <div className={`text-sm opacity-60 ${className}`} {...props}>
    {props.value}
  </div>
);

export type UiInputCautionProps = {
  value: string;
} & DivProps;
export const UiInputCaution = ({ className = '', ...props }: UiInputCautionProps) => (
  <div className={`text-sm text-my-red opacity-80 ${className}`} {...props}>
    {props.value}
  </div>
);
