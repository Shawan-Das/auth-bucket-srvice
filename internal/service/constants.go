package service

const ERP_ROUTE = "http://118.67.213.45:8081"
const COMPANY_LIST_API = ERP_ROUTE + "/api/external/getCompanyIds"
const UNIT_LIST_API = ERP_ROUTE + "/api/external/getUnitByCompany"
const ONE_UNIT_DATA = ERP_ROUTE + "/api/external/getUnitByUnitId"
const PAYROLL_API_BASE = "/api/v1/payroll"
const ADDL_PAYMENT_API_BASE = "/api/v1/addlpayments"
const INCOME_TAX_API_BASE = "/api/v1/incometax"

const AUTH_API_BASE = "/api/auth/login"
const UTIL_API_BASE = "/api/v1/utils"
const REF_API_BASE = "/api/v1/refdata"
const RPT_API_BASE = "/api/v1/report"
const HRM_API_BASE = "/api/v1/employee"

const SALARY_STATUS = "FINAL"
const GENERATE_SLARY_ACTION = "GENERATE-SALARY"
const CALCULATE_CPF_EMI_ACTION = "CALC-CPF-EMI"
const ALLOT_CPF_LOAN_ACTION = "ALLOT-CPF-LOAN"
const PAUSE_LOAN_ACTION = "PAUSE-LOAN"

const DEACTIVATE_LOAN_ACTION = "DEACTIVATE-LOAN"
const CRUD_ADDL_PAYMENT = "CRUD-ADDL-PAYMENT"
const APPROVE_REJECT_BONUS_INCENTIVE_PAYMENT = "APPROVE-REJCT-BONUS-INCENTIVE-PAYMENT"
const APPROVE_REJECT_OVERTIME_PAYMENT = "APPROVE-REJCT-OVERTIME-PAYMENT"
const CRUD_SALARY_STRUCTURE = "CRUD-SALARY-STRUCTURE"
const CRUD_PAYROLL_CONF = "CRUD-PAYROLL-CONF"
const PAYROLL_VIW = "PAYROLL-VIEW"
const BATCH_STATUS_SUCCESS = "SUCCESS"
const BATCH_STATUS_FAILED = "FAILED"
const CRU_ITAX_SCH3 = "CRU_ITAX_SCH3"
const APPROVE_REJECT_ITAX_SCH3 = "APPROVE_REJCT_CRU_ITAX_SCH3"

const LOAN_TYPE_CPF = 0
const LOAN_TYPE_SALARY_ADV = 1
const LOAN_STATUS_ACTIVE = "ACTIVE"
const LOAN_STATUS_CLOSED = "CLOSED"
const FESTIVAL_PAYMENT_TYPE = "FESTIVAL-BONUS-INCENTIVE"
const AREAR_PAYMENT_TYPE = "AREAR"
const OVERTIME_PAYMENT_TYPE = "OVERTIME"
const ADDITIONAL_PAYMENT_TYPE_ALLOWANCE = "ALLOWANCE"

const DATE_FORMAT = "2006/01/02"
const BONUS_SUBMITTED = "SUBMITTED"
const OT_SUBMITTED = "SUBMITTED"
const INVESTMENT_SUBMITTED = "SUBMITTED"
const INVESTMENT_APPROVED = "APPROVED"
const INVESTMENT_RETURNED = "RETURNED"

const BONUS_APPROVED = "APPROVED"
const BONUS_RETURNED = "RETURNED"
const ADDL_PAYMENT_PENDING = "PAYMENT_PENDING"
const ADDL_PAYMENT_PAID = "PAID"

const BONUS_PAID = "PAID"
const DRAFT = "DRAFT"
const AREAR_PENDING_APPROVAL = "PENDING_APPROVAL"
const AREAR_APPROVED = "APPROVED"
const AREAR_RETURNED = "RETURNED"
const OT_APPROVED = "APPROVED"
const OT_RETURNED = "RETURNED"
const OT_PAID = "PAID"
const DEPARTMENT_CRUD = "DEPARTMENT_CRUD"
const ADD_NEW_LEAVE_MASTER = "ADD_NEW_LEAVE_MASTER"
const EMP_EMPLOYMENT_INFO = "EMP_EMPLOYMENT_INFO"
const EMP_LEAVE_ENT = "EMP_LEAVE_ENT"
const APPLY_EMP_LEAVE = "APPLY_EMP_LEAVE"
const ADD_NEW_SHIFT = "ADD_NEW_SHIFT"
const DELETE_EMPLOYEE_ACCESS = "DELETE_EMPLOYEE_ACCESS"
const DIVISION_CRUD = "DIVISION_CRUD"
const ORG_HEAD_CRUD = "ORG_HEAD_CRUD"
const ASSIGN_ORG_HEAD = "ASSIGN_ORG_HEAD"

// Operation Table Constants
const ORG_HEAD = "ORG-HEAD"
const REFERENCE = "REFERENCE"
const EMPINFO_DELETE = "EMPLOYEE_GENERAL_INFO,EMPLOYEE_EMPLOYMENT_INFO,EMPLOYEE_IMAGES,EMPLOYEE_EXTN_INFO"
const EMP_EXTRA_INFO = "EMPLOYEE_EXTRA_INFO"
const EMP_GEN_INFO = "EMPLOYEE_GENERAL_INFO"

// General constants
const STATUS_ACTIVE = "ACTIVE"

// Cache constants
const EMP_SALARY_CACHE_KEY = "emp_saraly_info"
const ADDL_PAYMENT_TYPES_CACHE_KEY = "addl_payment_type"

// Default values
const DEFAULT_DATE = "2000/01/01"
const PAYROLL_CUTOFF_DATE = 15
const OVERTIME_CUTOFF_DATE = 25
const FileSizeInMB = 1024 * 1024

// Default password for new emp user
const DFAULT_PASS = "passw0rd"

// SHARED EMOPLOYEE QUERY
const SharedEmpListBaseQuery = `SELECT sl_id, emp_code, emp_identity_number, first_name, last_name, short_name,
emp_email, mobile_no, father_name, father_ocupation, mother_name, mother_ocupation, spouse_name, number_of_child,
nationality, weight, height, gender, date_of_birth, religion, meritalstatus, bloodgroup, status, cardid, suffix,
designation, company_id, company_name, unit_id, unit_name, division_id, division_name, dept_id, dept_name, section_id,
section_name, sub_section_id, sub_section_name, nationalid, birth_cert_no, tinno, joindate, passport_no, passport_issue_date,
passport_exp_date, visa_no, visa_issue_date, visa_expire_date, profile_picture, signature, educationinfo
FROM hrm.vw_employee_partial_info `
const getSharedEmployeeInfo = SharedEmpListBaseQuery + ` ORDER BY emp_code;`
const shareEmployeeOnSearch = SharedEmpListBaseQuery + `
WHERE (
    LOWER(full_name) LIKE '%'||LOWER($1)||'%' OR 
    (mobile_no) LIKE '%' || $1 || '%' OR
    emp_code = $1 OR
    emp_email = LOWER($1)
) ORDER BY emp_code;`
const shareEmployeeByCompany = SharedEmpListBaseQuery + `
    WHERE company_id = $1
    ORDER BY emp_code;`
const shareOneEmployeeData = SharedEmpListBaseQuery + ` WHERE emp_code = $1`

// const SharedEmpListBaseQuery = `
// SELECT G.id as sl_id, G.empcode AS emp_code, G.emp_identity_number, G.empname AS first_name, G.lastname AS last_name,
//     G.shortname AS short_name, G.email AS emp_email, (G.country_code || G.mobile_no) AS mobile_no,
//     G.fname AS father_name, G.fatheroccupation AS father_ocupation, G.mname AS mother_name,
//     G.motheroccupation AS mother_ocupation, G.spouse_name,  COALESCE(NULLIF(G.child, '')::int, 0) as number_of_child,
//     G.nationality, G.weight, G.height, G.gender, G.dob AS date_of_birth, G.religion, G.meritalstatus,
//     G.bloodgroup,
//     CASE
//         WHEN G.status = 'ACTIVE' THEN 1
//         ELSE 0
//     END AS status , G.cardid, G.suffix, R.diplaytext AS designation, G.companyid AS company_id,
//     G.company_name, G.unit_id, G.unit_name, G.division_id, G.division_name, G.dept_code AS dept_id, G.dept_name,
//     G.section_code AS section_id, G.section_name, G.sub_section_code AS sub_section_id, G.sub_section_name,
//     G.nationalid, G.birth_cert_no, G.tinno, EE.joindate, G.pasportno AS passport_no, G.passportissuedate AS passport_issue_date,
//     G.passportexpaireddate AS passport_exp_date,
//     G.visa_no, G.visa_issue_date, G.visa_expire_date,
//     CASE
//         WHEN i.profileimg LIKE 'https://satcomerp.s3.ap-south-1.amazonaws.com%' THEN I.profileimg
//         ELSE null
//     END AS profile_picture, I.signatureimg AS signature, eei.educationinfo
// FROM hrm.employee_general_info G
// JOIN hrm.employee_images I ON G.empcode = I.empcode
// JOIN hrm.employee_employment_info EE ON G.empcode = EE.empcode
// JOIN hrm.reference_info R ON EE.designation = R.value AND R.category = 'DESIGNATION'
// LEFT JOIN hrm.employee_extn_info eei on EEI.empcode = G.empcode
// `
// const getSharedEmployeeInfo = SharedEmpListBaseQuery + `
//     ORDER BY G.empcode`

// const shareEmployeeOnSearch = SharedEmpListBaseQuery + `
// WHERE (
//     LOWER(G.empname||' '||G.lastname) LIKE '%'||LOWER($1)||'%' OR
//     (country_code||mobile_no) LIKE '%' || $1 || '%' OR
//     G.empcode = $1 OR
//     g.email = LOWER($1)
// ) ORDER BY G.empcode`

// const shareEmployeeByCompany = SharedEmpListBaseQuery + `
//     WHERE G.companyid = $1
//     ORDER BY G.empcode`

// const shareOneEmployeeData = SharedEmpListBaseQuery + `
//     WHERE G.empcode = $1`
