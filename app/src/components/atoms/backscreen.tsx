import React from 'react';
import { DivProps } from './type';

export type BackscreenProps = DivProps;

export const Backscreen = ({ className = '', ...props }: BackscreenProps) => (
  <div className={`bg-my-black text-my-white ${className}`} {...props} />
);
