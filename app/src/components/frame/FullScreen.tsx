import { FadeAnim, FadeAnimProps } from 'components/anims/FadeAnim';
import { Backscreen } from 'components/atoms/backscreen';
import React from 'react';

export type FullScreenTitleProps = {
  active?: boolean;
} & Pick<FadeAnimProps, 'children'>;
export const FullScreentitle = ({ active }: FullScreenTitleProps) => {
  return (
    <Backscreen className='absolute'>
      <FadeAnim in={Boolean(active)}>
        <div className='flex h-screen w-screen items-center justify-center'>
          <h1 className='text-3xl'>
            Ctrl<span className='text-my-yellow'>+</span>
          </h1>
        </div>
      </FadeAnim>
    </Backscreen>
  );
};
