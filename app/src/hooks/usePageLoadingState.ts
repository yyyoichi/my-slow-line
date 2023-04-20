import { useLoadingState } from './useLoadingState';
import { usePageState } from './usePageState';

/**manage page send and load. */
export const usePageLoadingState = <T>(defaultPage: T, defaultLoad: boolean) => {
  const { currentPage, setCurrentPage } = usePageState<T>(defaultPage);
  const { isLoading, setIsLoading } = useLoadingState(defaultLoad);

  const _goToPageAndReleaseLoad = (page: T) => {
    setCurrentPage(page);
    setIsLoading(false);
  };
  const presetPageAndStartLoading = (newPage: T) => {
    const priviousPage = currentPage;
    setIsLoading(true);
    const goToNextPage = () => _goToPageAndReleaseLoad(newPage);
    const resetCurrentPage = () => _goToPageAndReleaseLoad(priviousPage);
    return {
      goToNextPage,
      resetCurrentPage,
    };
  };

  return {
    currentPage,
    isLoading,
    presetPageAndStartLoading,
  };
};
