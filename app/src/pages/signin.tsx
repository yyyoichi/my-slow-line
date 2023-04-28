import dynamic from 'next/dynamic';

export default dynamic(() => import('templates/signin/'), { ssr: false });
