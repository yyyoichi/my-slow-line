import React from 'react';

/**manage page sending. */
export const usePageState = <T>(defaultPage: T) => {
  const [currentPage, setCurrentPage] = React.useState<T>(defaultPage);

  return {
    currentPage,
    setCurrentPage,
  };
};
