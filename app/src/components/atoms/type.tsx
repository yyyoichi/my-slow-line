import React from 'react';

type WithClassName = {
  className?: string;
};

export type DivProps = React.ComponentProps<'div'> & WithClassName;
export type InputProps = React.ComponentProps<'input'> & WithClassName;

export type MyColorPallet = 'red' | 'yellow' | 'light-black' | 'black' | 'white' | 'green';
