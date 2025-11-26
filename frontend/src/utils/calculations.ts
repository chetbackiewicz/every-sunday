import type { CategoryNode, Transaction, CellReference } from '../types';

// Helper to flatten tree for easy lookup
export const flattenCategories = (nodes: CategoryNode[]): CategoryNode[] => {
  let flat: CategoryNode[] = [];
  for (const node of nodes) {
    flat.push(node);
    if (node.children) {
      flat = flat.concat(flattenCategories(node.children));
    }
  }
  return flat;
};

export const calculateCategoryActuals = (
  categories: CategoryNode[],
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  _transactions: Transaction[],
  cellReferences: CellReference[]
): CategoryNode[] => {
  // 1. Create a map of category ID -> sum of linked transactions
  const linkSums = new Map<number, number>();
  
  for (const ref of cellReferences) {
     const current = linkSums.get(ref.categoryId) || 0;
     linkSums.set(ref.categoryId, current + ref.amount);
  }

  // 2. Recursive function to calculate actuals (including children)
  const processNode = (node: CategoryNode): CategoryNode => {
     // Calculate self actual from links
     const selfActual = linkSums.get(node.id) || 0;

     // Process children
     let childrenSum = 0;
     const processedChildren = node.children?.map(child => {
        const processed = processNode(child);
        // Use manual actual if present, otherwise calculated
        const effectiveActual = processed.manualActual ?? processed.calculatedActual ?? 0;
        childrenSum += effectiveActual;
        return processed;
     });

     // Total calculated actual = self links + children effective actuals
     const calculatedActual = selfActual + childrenSum;

     return {
        ...node,
        calculatedActual,
        children: processedChildren
     };
  };

  return categories.map(processNode);
};

export interface BudgetSummary {
  income: number;
  expenses: number;
  savings: number; // Pre + Post tax
  remaining: number;
}

export const calculateBudgetSummary = (categories: CategoryNode[]): BudgetSummary => {
   // This assumes top-level categories have specific names or types
   // For this MVP, let's assume:
   // - Income
   // - Expenses
   // - Savings
   
   // We need a way to identify them. For now, we'll match by name loosely or ID if fixed.
   // Let's flatten to find them.
   
   const flat = flattenCategories(categories);
   
   const getActual = (name: string) => {
       const cat = flat.find(c => c.name.toLowerCase() === name.toLowerCase());
       if (!cat) return 0;
       return cat.manualActual ?? cat.calculatedActual ?? 0;
   };

   const income = getActual('Income');
   const expenses = getActual('Expenses');
   // Assuming Savings is a top level or tracked
   const savings = getActual('Savings'); // Or sum of Pre-Tax / Post-Tax if they exist

   return {
      income,
      expenses,
      savings,
      remaining: income - expenses - savings
   };
};

