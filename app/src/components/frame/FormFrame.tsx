import React from 'react';
// wrap contents
const Container = ({ children }: { children: React.ReactNode }) => (
  <div className='flex flex-col gap-3 rounded-md bg-my-light-black py-2 '>{children}</div>
);
// wrap section in contents
const Content = ({ children }: { children: React.ReactNode }) => <div className='px-4 py-1'>{children}</div>;

const c = {
  Container,
  Content,
};

export default c;
