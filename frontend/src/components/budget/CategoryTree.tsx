import React, { useState } from 'react';
import { Tree } from 'react-arborist';
import { useAppState } from '../../hooks/useAppState';
import { useBudgetCalculations } from '../../hooks/useBudgetCalculations';

// Simple node component with cell linking support
const Node = ({ node, style, dragHandle }: any) => {
  const { state, dispatch } = useAppState();
  const isLinked = state.activeCategoryId === node.id && state.isLinkingMode;

  const handleClick = () => {
      if (state.isLinkingMode && state.activeCategoryId === node.id) {
           dispatch({ type: 'SET_LINKING_MODE', payload: { active: false, categoryId: null } });
      } else {
           dispatch({ type: 'SET_LINKING_MODE', payload: { active: true, categoryId: node.id } });
      }
  };

  const handleAdd = (e: React.MouseEvent) => {
    e.stopPropagation();
    // Placeholder for Add logic - ideally opens a modal or adds a temp node
    console.log('Add child to', node.id);
    // Dispatch ADD_CATEGORY action (mock)
    // dispatch({ type: 'ADD_CATEGORY', payload: { parentId: node.id, name: 'New Category', projected: 0 } });
  };

  const handleDelete = (e: React.MouseEvent) => {
    e.stopPropagation();
    if (confirm(`Delete ${node.data.name}?`)) {
       console.log('Delete', node.id);
       // Dispatch DELETE_CATEGORY action (mock)
    }
  };

  return (
    <div 
      style={style} 
      ref={dragHandle} 
      className={`flex items-center py-1 cursor-pointer group ${isLinked ? 'bg-indigo-100 ring-2 ring-indigo-500' : 'hover:bg-gray-50'}`}
      onClick={handleClick}
    >
       <span className="mr-2 text-gray-400">
         {node.isLeaf ? '•' : (node.isOpen ? '▼' : '▶')}
       </span>
       <span className="flex-1 truncate text-sm text-gray-700">{node.data.name}</span>
       <span className="text-sm text-gray-500 mr-2">
          {node.data.manualActual !== null && node.data.manualActual !== undefined ? '✏️ ' : ''} 
          ${(node.data.manualActual ?? node.data.calculatedActual ?? 0).toFixed(2)} / ${node.data.projected.toFixed(2)}
       </span>
       
       <div className="hidden group-hover:flex space-x-1 mr-1">
          <button onClick={handleAdd} className="text-xs text-blue-500 hover:text-blue-700 px-1" title="Add Child">+</button>
          <button onClick={handleDelete} className="text-xs text-red-500 hover:text-red-700 px-1" title="Delete">×</button>
       </div>
    </div>
  );
};

export const CategoryTree: React.FC = () => {
  const { state } = useAppState();
  const { enrichedCategories } = useBudgetCalculations();
  const [term, setTerm] = useState('');

  // Use enriched categories if available, else state (or mock)
  // Note: enrichedCategories will be computed from state.categories
  // If state.categories is empty, we might want mock data for demo
  const data = state.categories.length > 0 ? enrichedCategories : [
    { id: '1', name: 'Income', projected: 5000, calculatedActual: 0, manualActual: null, children: [] },
    { id: '2', name: 'Expenses', projected: 3000, calculatedActual: 0, manualActual: null, children: [
        { id: '3', name: 'Housing', projected: 1500, calculatedActual: 0, manualActual: null, children: [] },
        { id: '4', name: 'Food', projected: 500, calculatedActual: 0, manualActual: null, children: [] }
    ]}
  ];

  return (
    <div className="h-full min-h-[400px]">
      <div className="mb-2">
         <input 
            type="text" 
            placeholder="Search categories..." 
            className="w-full border-gray-300 rounded-md shadow-sm sm:text-sm p-2 border"
            value={term}
            onChange={(e) => setTerm(e.target.value)}
         />
      </div>
      <Tree
        initialData={data as any} 
        searchTerm={term}
        openByDefault={true}
        width={300}
        height={600}
        indent={24}
        rowHeight={32}
        overscanCount={1}
        paddingTop={30}
        paddingBottom={10}
        padding={25}
      >
        {Node}
      </Tree>
    </div>
  );
};
