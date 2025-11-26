import { useMemo } from 'react';
import { useAppState } from './useAppState';
import { calculateCategoryActuals, calculateBudgetSummary } from '../utils/calculations';

export const useBudgetCalculations = () => {
  const { state } = useAppState();

  const enrichedCategories = useMemo(() => {
    return calculateCategoryActuals(state.categories, state.transactions, state.cellReferences);
  }, [state.categories, state.transactions, state.cellReferences]);

  const summary = useMemo(() => {
    return calculateBudgetSummary(enrichedCategories);
  }, [enrichedCategories]);

  return {
    enrichedCategories,
    summary
  };
};

