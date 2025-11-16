# Budget Tracking Application - Research Documentation

This directory contains detailed research findings for the technical spike investigation.

## Status Overview

**Frontend Research**: 80% Complete (4/5 done)
- ✅ CSV Parsing & File Upload
- ✅ Editable Tables/Grids
- ✅ Hierarchical Category Tree
- ✅ Cell Linking/Reference System
- ⏳ Real-time Calculations (pending)

**Backend Research**: 0% Complete (0/2 done)
- ⏳ Golang Framework
- ⏳ PostgreSQL Schema

---

## Research Documents

### ✅ Completed Research

1. **[01-csv-parsing.md](./01-csv-parsing.md)** - CSV Parsing & File Upload
   - **Decision**: PapaParse v5 + react-dropzone
   - Streaming support, 10-file limit, validation

2. **[02-editable-tables.md](./02-editable-tables.md)** - Editable Table/Grid Components
   - **Decision**: TanStack Table (React Table v8)
   - Headless architecture, lightweight (10-15kb), custom cell editing

3. **[03-category-tree.md](./03-category-tree.md)** - Hierarchical Category Tree UI
   - **Decision**: react-arborist
   - Depth control, drag & drop, virtual scrolling

4. **[04-cell-linking.md](./04-cell-linking.md)** - Cell Linking/Reference System
   - **Decision**: useReducer with Lifted State
   - Complex state coordination, referential integrity

### ⏳ Pending Research

5. **Real-time Calculation Patterns** - Not started
   - Auto-calculation for Projected/Actual/Difference
   - Remaining formula (Income - Savings - Expenses)
   - Cascading updates, memoization

6. **Multi-month Data Management** - Not started
   - Month switching pattern
   - Data persistence per month
   - URL routing vs dropdown

7. **Golang Backend Architecture** - Not started
   - Framework selection (Gin/Echo/Fiber)
   - REST API design
   - File upload handling

8. **PostgreSQL Schema Design** - Not started
   - Hierarchical categories modeling
   - Transaction tracking
   - Cell references storage

---

## Technology Stack Summary

### Confirmed Choices

**Frontend**:
- CSV Parsing: PapaParse v5
- File Upload: react-dropzone
- Tables: TanStack Table v8
- Category Tree: react-arborist
- State Management: useReducer at app root

**Backend** (pending):
- Framework: TBD (Gin/Echo/Fiber)
- Database: PostgreSQL
- ORM: TBD

---

## Reading Order

For understanding the complete architecture:

1. Start with `01-csv-parsing.md` - data ingestion
2. Then `02-editable-tables.md` - data display/editing
3. Then `03-category-tree.md` - budget structure
4. Then `04-cell-linking.md` - data relationships
5. Finally `spike-doc.md` - overall progress and next steps

---

## Contributing to Research

When adding new research:
1. Create a new numbered markdown file (e.g., `05-calculations.md`)
2. Follow the existing template structure:
   - Status, Date, Decision at top
   - Recommendation Summary
   - Detailed findings with code examples
   - Pros/Cons
   - Comparison tables if applicable
3. Update this README with the new document
4. Update `spike-doc.md` with decision summary
