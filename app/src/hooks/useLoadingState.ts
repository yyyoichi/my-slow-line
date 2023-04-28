import React from 'react';

export const useLoadingState = (defaultLoad: boolean) => {
  const [isLoading, setIsLoading] = React.useState<boolean>(defaultLoad);
  return {
    isLoading,
    setIsLoading,
  };
};
