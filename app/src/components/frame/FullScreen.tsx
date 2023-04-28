import { FadeAnim, FadeAnimProps } from 'components/anims/FadeAnim';
import { Backscreen } from 'components/atoms/backscreen';
import React from 'react';

export type FullScreenTitleProps = {
  active?: boolean;
} & Pick<FadeAnimProps, 'children'>;
export const FullScreentitle = ({ active }: FullScreenTitleProps) => {
  return (
    <Backscreen className='absolute overflow-hidden'>
      <FadeAnim in={Boolean(active)}>
        <h1 className='flex min-h-screen w-screen items-center justify-center text-3xl'>
          Ctrl<span className='text-my-yellow'>+</span>
        </h1>
      </FadeAnim>
    </Backscreen>
  );
};
