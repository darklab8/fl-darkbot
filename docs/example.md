# Mkdocs Examples

## Requirements

- [Head First Git](https://www.amazon.com/Head-First-Git-Learners-Understanding/dp/1492092517)
- [Learning GIT CLI intectively](https://learngitbranching.js.org/)
- [Learning advanced Git](https://git-scm.com/book/en/v2)

<!DOCTYPE html>

<html lang="en">
<head>
  <meta charset="utf-8">
</head>
<body>
  <div class="mermaid">
flowchart TD
  GitCommit[attempt to fixate commit like\ngit commit -m 'feat: add rendering in markdown format'\nwith 'autogit hook activate' enabled]
  RequestValidatingChangelog[Request changelog with --validate flag] --> TryParsingCommitMessage
  GitCommit --> TryParsingCommitMessage[Try parsing commit message\nto git conventional commit\ntype \ scope \ subject \ body \ footers]
  TryParsingCommitMessage --> ReportFail[Reporting errors if unable]
  TryParsingCommitMessage --> ContinuingValidation[Continue Validation]
  ContinuingValidation --> CheckOptionalValidationRulesIfEnabled[Check options validation rules\nif they are enabled]
  CheckOptionalValidationRulesIfEnabled --> CommitTypeInAllowedList[Commit type is\nin allowed list]
  CommitTypeInAllowedList --> WhenAppliedRules
  CheckOptionalValidationRulesIfEnabled --> MinimumNWords[Minimum N words is present\nin commit subhect]
  MinimumNWords --> WhenAppliedRules
  CheckOptionalValidationRulesIfEnabled --> IssueIsLinked[Issue is linked to commit]
  IssueIsLinked --> WhenAppliedRules
  CheckOptionalValidationRulesIfEnabled --> CheckOtherEnabledRulesInSettings[Check other enabled\nrules in settings]
  CheckOtherEnabledRulesInSettings --> WhenAppliedRules
  CheckOptionalValidationRulesIfEnabled --> WhenAppliedRules[when applied rules]
  WhenAppliedRules --> IfCommit[if it was commit,\nthen fixate if passed rules,\nor cancel fixation]
  WhenAppliedRules --> IfChangelog[if it was changelog validation\nthen report no errors and exit code 0\nfor pipeline checks]
  </div>
 <script src="mermaid.min.js"></script>
 <script src="/mermaid.min.js"></script>
 <script>mermaid.initialize({startOnLoad:true});
</script>
</body>
</html>

### Contacts

- [Darklab Discord server](https://discord.gg/aukHmTK82J)
- [feature requests can be opened in Github Issues](https://github.com/darklab8/fl-darkbot/issues)
- email to `dark.dreamflyer@gmail.com`
