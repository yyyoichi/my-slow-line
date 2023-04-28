import React from 'react';
import { DivProps, MyColorPallet } from './type';

export type UiButtonProps = {
  color: MyColorPallet;
} & DivProps;
export const UiButton = ({ className = '', color, ...props }: UiButtonProps) => {
  const buttonProps: DivProps = {
    ...props,
  };
  return (
    <div className={`rounded-md text-center font-bold ${getColorClassName(color)} ${className}`} {...buttonProps} />
  );
};

const getColorClassName = (color: UiButtonProps['color']) => {
  switch (color) {
    case 'red':
      return 'bg-my-red text-my-white';
    case 'yellow':
      return 'bg-my-yellow text-my-black';
    case 'light-black':
      return 'bg-my-light-black text-my-white';
    case 'black':
      return 'bg-my-black text-my-white';
    case 'white':
      return 'bg-my-white text-my-black';
    case 'green':
      return 'bg-my-green text-my-black';
  }
};
