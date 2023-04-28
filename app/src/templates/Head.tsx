import { UiHead } from 'components/frame';
import { MyAccountContext } from 'hooks';
import Link from 'next/link';
import React from 'react';

export const Head = () => {
  const ac = React.useContext(MyAccountContext);
  return <UiHead>{ac.myAccount.has ? <User name={ac.myAccount.name} /> : <Guest />}</UiHead>;
};

const User = ({ name }: { name: string }) => {
  return (
    <div className='ml-auto mr-4 flex aspect-square w-10 items-center justify-center rounded-full bg-my-white p-1 pb-2 text-xl font-bold text-my-black'>
      {name[0]}
    </div>
  );
};

const Guest = () => {
  return (
    <div className='ml-auto px-4'>
      <Link prefetch={true} className='mx-2 inline-block rounded border-2 border-my-white px-3 py-1' href={'/login'}>
        login
      </Link>
      <Link
        prefetch={true}
        className='mx-2 inline-block rounded border-2 border-my-white bg-my-white px-3 py-1 font-semibold text-my-black'
        href={'/signin'}
      >
        signin
      </Link>
    </div>
  );
};
