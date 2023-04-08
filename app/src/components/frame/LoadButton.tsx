import { UiButton, UiButtonProps } from 'components/atoms/button';
import { UiLoader } from 'components/atoms/loader';
import { MyColorPallet } from 'components/atoms/type';
import React from 'react';

export type LoadButtonProps = {
  active: boolean;
} & UiButtonProps;

export const LoadButton = ({ active, className = '', color, children, onClick, ...props }: LoadButtonProps) => {
  const Child = active ? children : <UiLoader color={getPairColor(color)} />;
  const onclick = active ? onClick : undefined;
  return (
    <UiButton
      className={`w-full text-center ${active ? 'cursor-pointer' : ''} ${className}`}
      color={color}
      onClick={onclick}
      {...props}
    >
      {Child}
    </UiButton>
  );
};
/**get vilibility color to [color] */
const getPairColor = (color: MyColorPallet): MyColorPallet => {
  switch (color) {
    case 'red':
      return 'white';
    case 'yellow':
      return 'black';
    case 'light-black':
      return 'white';
    case 'black':
      return 'white';
    case 'white':
      return 'black';
    case 'green':
      return 'black';
  }
};
