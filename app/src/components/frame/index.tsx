import React from 'react';
import { DivProps } from 'components/atoms/type';

export const UiMain = ({ className = '', ...props }: DivProps) => (
  <div className='bg-my-black'>
    <div className={`flex min-h-screen w-full flex-col text-my-white md:mx-auto md:w-2/3 ${className}`} {...props}>
      {props.children}
    </div>
  </div>
);

export const UiHead = (props: React.ComponentProps<'div'>) => (
  <div className={`flex h-10 w-full`} {...props}>
    <div className='px-3 py-1 text-2xl'>Ctrl+</div>
  </div>
);

export const UiYCenter = ({ className = '', ...props }: DivProps) => (
  <div className={`my-auto ${className}`}>{props.children}</div>
);
