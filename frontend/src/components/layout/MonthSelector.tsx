import React from 'react';

export const MonthSelector: React.FC = () => {
  return (
    <div className="mb-4">
      <label htmlFor="month" className="block text-sm font-medium text-gray-700">
        Select Month
      </label>
      <input
        type="month"
        id="month"
        name="month"
        className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
      />
    </div>
  );
};

