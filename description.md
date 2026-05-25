Use this tool when you need to ask the user questions during execution. This allows you to:
1. Gather user preferences or requirements
2. Clarify ambiguous instructions
3. Get decisions on implementation choices as you work
4. Offer choices to the user about what direction to take.

Usage notes:
- Users will always be able to select "Other" to provide custom text input
- Use multiSelect: true to allow multiple answers to be selected for a question
- If you recommend a specific option, make that the first option in the list and add "(Recommended)" at the end of the label
- When the questions belong to a specific project, repo, or task, populate the optional `context` field with a short hint (e.g. "dw — drafting brief for idea 2026-05-25-foo", "primer-ui PR #42", "/Users/me/code/widget"). The host renders it above the questions so the user can see what work the questions are part of. The hint is shown to the user only; it never appears in the answer returned to you.
