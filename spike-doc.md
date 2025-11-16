Goal:
I'm building a web application that will allow users to track their money across different months and weeks
Create a web application that will allow the user to drag and drop a CSV file, example ../../CC_Statement_Example.CSV
The the application should display the different rows given from the CSV
- The category should be left blank however

For each month the user needs to fill out different areas with each have a "Projected," "Actual," "Difference"

Categories:
- Income
- Savings (Pre-Tax)
- Savings (Post-Tax)
- Expenses
- Remaining (Income - Savings - Expenses = )

For each category the user should be able to create subcategories and subcategories to those (limit the branch length to two nodes besides the parent node, like "Income")
- For example, in income the user should see a + sign, where they can add "Salary," "Commission," or "Dividends"
- In expenses, a user should see an + button where they could add "House," and underneath house should see another + button where they could add "Mortgage," "Insurance," "HOA Fees," etc.  They shouldn't see any + button under Insurance or Mortgage or HOA Fees however.

To the right of this monthly budget should be the drag and drop area where the user can drop in the bank statement's csv
 the user should be able to input up to ten CSV's
 when the user drops the CSV in, a table should populate with the Transaction Date,Post Date OR Posting Date,Description,Category,Type,Amount, and Note
 - Each of the cells should be editable
 - The Category and Note should be left blank

 To assign any expense or credit to a category, the user will need to click the cell next to each category, under the "Actual" column. This will color the cell
 Then the user needs to go and click on the cells from the uploaded csv files


 Each cell should be able to be deleted with a red "X" next to the row
 - If a cell is deleted and is referenced in any of the category "Actual" cells, the deleted value should be removed from them as well


The values of the cells of any area should be able to be overriden by the user clicking on the cell and typing, even the "Actual" categories

The calculations of the total income, savings, and other categories should be automatically updated as values are added, either manually when filling out the categories like "Income" or "Expenses" or when the users has clicked on an "Actual" cell and has started add values to it.

The user should have different months in a drop down that they can create a budget tracking sheet like the above.


 

# Technology used

## Front end

- React + Typescript

## Backend

- Golang

## Database

- Postgresql


---

# Spike Investigation - Budget Tracking Application

**Investigation Start**: October 19, 2025  
**Status**: In Progress - Frontend 80% Complete

> **📁 Detailed Research**: See [`/research`](./research/) directory for complete findings

---

## Quick Reference - Technology Decisions

### ✅ Confirmed Choices

| Area                 | Decision               | Status     | Doc                                                       |
| -------------------- | ---------------------- | ---------- | --------------------------------------------------------- |
| **CSV Parsing**      | PapaParse v5           | ✅ Complete | [01-csv-parsing.md](./research/01-csv-parsing.md)         |
| **File Upload**      | react-dropzone         | ✅ Complete | [01-csv-parsing.md](./research/01-csv-parsing.md)         |
| **Editable Tables**  | TanStack Table v8      | ✅ Complete | [02-editable-tables.md](./research/02-editable-tables.md) |
| **Category Tree**    | react-arborist         | ✅ Complete | [03-category-tree.md](./research/03-category-tree.md)     |
| **State Management** | useReducer at App Root | ✅ Complete | [04-cell-linking.md](./research/04-cell-linking.md)       |
| **Cell Linking**     | Lifted State Pattern   | ✅ Complete | [04-cell-linking.md](./research/04-cell-linking.md)       |

### ⏳ Pending Research

| Area                       | Status      | Priority | Notes                                                        |
| -------------------------- | ----------- | -------- | ------------------------------------------------------------ |
| **Real-time Calculations** | Not Started | HIGH     | Auto-calc for Projected/Actual/Difference, cascading updates |
| **Multi-month Management** | Not Started | MEDIUM   | Month switching, data persistence                            |
| **Golang Backend**         | Not Started | MEDIUM   | Framework selection (Gin/Echo/Fiber)                         |
| **PostgreSQL Schema**      | Not Started | MEDIUM   | Hierarchical categories, transactions, cell references       |

---

## Research Questions

### Frontend (React + TypeScript)
1. ✅ **CSV File Handling** - COMPLETE → PapaParse v5 + react-dropzone
2. ✅ **Editable Table Component** - COMPLETE → TanStack Table v8
3. ✅ **Hierarchical Category UI** - COMPLETE → react-arborist
4. ✅ **Cell Linking System** - COMPLETE → useReducer with lifted state
5. ⏳ **Real-time Calculations** - PENDING
6. ✅ **State Management** - COMPLETE → useReducer at app root

### Backend (Golang)
7. ⏳ **API Framework** - PENDING
8. ⏳ **Database ORM** - PENDING

### Database (PostgreSQL)
9. ⏳ **Schema Design** - PENDING
10. ⏳ **Reference Tracking** - PENDING (partially addressed in cell linking research)
11. ⏳ **Multi-month Support** - PENDING

## Core Requirements Summary

Based on initial analysis (Oct 19, 2025):

- **Multi-file CSV upload**: Support up to 10 CSV files per month
- **Editable grid**: All cells in CSV table must be editable (Transaction Date, Post Date, Description, Category, Type, Amount, Note)
- **Category hierarchy**: 3-level max depth (Parent → Child → Grandchild), e.g., Expenses → House → Mortgage
- **Cell linking interaction**: Click "Actual" cell to activate → click CSV cells to add values
- **Referential integrity**: Deleting CSV row removes values from linked "Actual" cells
- **Auto-calculations**: Income - Savings (Pre-Tax) - Savings (Post-Tax) - Expenses = Remaining
- **Month selector**: Dropdown to switch between different monthly budgets

---

## Research Progress

### Completed Investigations

Detailed findings available in `/research` directory:

1. **CSV Parsing & File Upload** → [01-csv-parsing.md](./research/01-csv-parsing.md)
2. **Editable Tables/Grids** → [02-editable-tables.md](./research/02-editable-tables.md)
3. **Hierarchical Category Tree** → [03-category-tree.md](./research/03-category-tree.md)
4. **Cell Linking/Reference System** → [04-cell-linking.md](./research/04-cell-linking.md)

### Active Investigation

#### 5. Real-time Calculation Patterns (Next Priority)

**Research Questions**:
- How to implement auto-calculation for Projected/Actual/Difference columns?
- How to calculate Remaining formula: Income - Savings (Pre-Tax) - Savings (Post-Tax) - Expenses?
- What React patterns work best for cascading updates?
- How to optimize performance with memoization?
- How to handle circular dependencies in calculations?

**Investigation Areas**:
- useMemo for expensive calculations
- useEffect vs reducer for calculation triggers
- Derived state patterns
- Performance profiling with React DevTools
- Debouncing for user input

---

## Archived Research (Completed)

> **Note**: Full details moved to `/research` directory. Summaries below for quick reference.

### CSV Parsing & File Upload Research

**Decision**: PapaParse v5 + react-dropzone  
**Key Features**: Streaming support (10MB chunks), 10-file limit via `maxFiles`, header detection  
**Full Details**: [01-csv-parsing.md](./research/01-csv-parsing.md)

### Editable Table/Grid Component Research

**Decision**: TanStack Table (React Table v8)  
**Why**: Headless architecture (10-15kb), complete design control, free and open source  
**Rejected**: AG Grid (500kb bundle, paid license for features)  
**Full Details**: [02-editable-tables.md](./research/02-editable-tables.md)

### Hierarchical Category Tree UI Research

**Decision**: react-arborist  
**Why**: Headless architecture, built-in drag & drop, easy depth limiting via `node.level`  
**Rejected**: MUI X Tree View (paid features), react-beautiful-dnd (archived)  
**Full Details**: [03-category-tree.md](./research/03-category-tree.md)

### Cell Linking/Reference System Research

**Decision**: useReducer with Lifted State at App Root  
**Why**: Complex state coordination, atomic updates, referential integrity  
**Pattern**: Single reducer manages categories + CSV data + linking state  
**Full Details**: [04-cell-linking.md](./research/04-cell-linking.md)

---

## Next Steps

1. **Complete real-time calculations research** (HIGH priority)
   - Auto-calc patterns for Projected/Actual/Difference
   - Cascading updates for Remaining formula
   - Memoization strategy

2. **Multi-month data management** (MEDIUM priority)
   - Month switching UX
   - Data persistence strategy

3. **Backend research** (MEDIUM priority)
   - Golang framework selection
   - PostgreSQL schema design

4. **Begin implementation** (After frontend research complete)
   - Set up project structure
   - Implement chosen libraries
   - Build proof-of-concept

---

## Implementation Roadmap (Draft)

### Phase 1: Core Frontend (Weeks 1-2)
- Set up React + TypeScript project
- Implement CSV upload with PapaParse + react-dropzone
- Build editable table with TanStack Table
- Create category tree with react-arborist

### Phase 2: State & Interactions (Weeks 3-4)
- Implement useReducer state management
- Build cell linking system
- Add auto-calculations
- Multi-month management

### Phase 3: Backend (Weeks 5-6)
- Set up Golang API
- Design PostgreSQL schema
- Implement REST endpoints
- Add authentication

### Phase 4: Integration & Polish (Weeks 7-8)
- Connect frontend to backend
- Add data persistence
- Performance optimization
- User testing and refinement

---

## References

- [React Documentation](https://react.dev)
- [TanStack Table Docs](https://tanstack.com/table/latest)
- [react-arborist GitHub](https://github.com/brimdata/react-arborist)
- [PapaParse Docs](https://www.papaparse.com/docs)

---

**Last Updated**: October 19, 2025  
**Next Review**: After completing calculations research