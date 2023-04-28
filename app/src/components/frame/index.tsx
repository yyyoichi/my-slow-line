import React from 'react';
import { DivProps } from 'components/atoms/type';

export const UiMain = ({ className = '', ...props }: DivProps) => (
  <div className='bg-my-black'>
    <div className={`flex min-h-screen w-full flex-col text-my-white md:mx-auto md:w-2/3 ${className}`} {...props}>
      {props.children}
    </div>
  </div>
);

export type UiHeadProps = DivProps;
export const UiHead = ({ children, className = '', ...props }: UiHeadProps) => (
  <div className={`flex w-full items-center ${className}`} {...props}>
    <div className='px-5 py-1 text-4xl'>
      Ctrl<span className='text-5xl text-my-yellow'>+</span>
    </div>
    {children}
  </div>
);

export const UiYCenter = ({ className = '', ...props }: DivProps) => (
  <div className={`my-auto ${className}`}>{props.children}</div>
);
