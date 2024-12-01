import React, { useState, useEffect } from 'react';
import { Pagination, Modal, Select, notification } from 'antd'; // 导入 Pagination 和 Modal 组件
import './SearchComponent.css';

// 假设我们有以下省份和部门的数据结构
const provincesData = [
    {id: 0, name: 'All'},
    {id: 35, name: '中央'},
    {id: 1, name: '北京市'},
    {id: 2, name: '天津市'},
    {id: 3, name: '河北省'},
    {id: 4, name: '山西省'},
    {id: 5, name: '内蒙古自治区'},
    {id: 6, name: '辽宁省'},
    {id: 7, name: '吉林省'},
    {id: 8, name: '黑龙江省'},
    {id: 9, name: '上海市'},
    {id: 10, name: '江苏省'},
    {id: 11, name: '浙江省'},
    {id: 12, name: '安徽省'},
    {id: 13, name: '福建省'},
    {id: 14, name: '江西省'},
    {id: 15, name: '山东省'},
    {id: 16, name: '河南省'},
    {id: 17, name: '湖北省'},
    {id: 18, name: '湖南省'},
    {id: 19, name: '广东省'},
    {id: 20, name: '广西壮族自治区'},
    {id: 21, name: '海南省'},
    {id: 22, name: '重庆市'},
    {id: 23, name: '四川省'},
    {id: 24, name: '贵州省'},
    {id: 25, name: '云南省'},
    {id: 26, name: '西藏自治区'},
    {id: 27, name: '陕西省'},
    {id: 28, name: '甘肃省'},
    {id: 29, name: '青海省'},
    {id: 30, name: '宁夏回族自治区'},
    {id: 31, name: '新疆维吾尔自治区'},
    {id: 32, name: '台湾省'},
    {id: 33, name: '香港特别行政区'},
    {id: 34, name: '澳门特别行政区'}
];

const departmentsData = [
    {id: 0, name: 'All'},
    {id:1,name:'科技部'},
    {id:2,name:'教育部'},
    {id:3,name:'工信部'},
    {id:4,name:'国务院文件'},
    {id:5,name:'外交部'},
    {id:6,name:'发展和改革委员会'},
    {id:7,name:'民族事务委员会'},
    {id:8,name:'公安部'},
    {id:9,name:'安全部'},
    {id:10,name:'民政部'},
    {id:11,name:'司法部'},
    {id:12,name:'财政部'},
    {id:13,name:'人力资源和社会保障部'},
    {id:14,name:'自然资源部'},
    {id:15,name:'生态环境部'},
    {id:16,name:'住房和城乡建设部'},
    {id:17,name:'交通运输部'},
    {id:18,name:'水利部'},
    {id:19,name:'农业农村部'},
    {id:20,name:'商务部'},
    {id:21,name:'文化和旅游部'},
    {id:22,name:'卫生健康委员会'},
    {id:23,name:'退役军人事务部'},
    {id:24,name:'应急管理部'},
    {id:25,name:'人民银行'},
    {id:26,name:'审计署'},
    {id:27,name:'国有资产监督管理委员会'},
    {id:28,name:'海关总署'},
    {id:29,name:'税务总局'},
    {id:30,name:'市场监督管理总局'},
    {id:31,name:'金融监督管理总局'},
    {id:32,name:'广播电视总局'},
    {id:33,name:'体育总局'},
    {id:34,name:'统计局'},
    {id:35,name:'国际发展合作署'},
    {id:36,name:'医疗保障局'},
    {id:37,name:'机关事务管理局'},
    {id:38,name:'标准化管理委员会'},
    {id:39,name:'新闻出版署'},
    {id:40,name:'版权局'},
    {id:41,name:'互联网信息办公室'},
    {id:42,name:'中国科学院'},
    {id:43,name:'中国社会科学院'},
    {id:44,name:'中国工程院'},
    {id:45,name:'中国气象局'},
    {id:46,name:'中国银行保险监督管理委员会'},
    {id:47,name:'中国证券监督管理委员会'},
    {id:48,name:'信访局'},
    {id:49,name:'粮食和物资储备局'},
    {id:50,name:'能源局'},
    {id:51,name:'数据局'},
    {id:52,name:'国防科技工业局'},
    {id:53,name:'烟草专卖局'},
    {id:54,name:'移民管理局'},
    {id:55,name:'林业和草原局'},
    {id:56,name:'铁路局'},
    {id:57,name:'中国民用航空局'},
    {id:58,name:'邮政局'},
    {id:59,name:'文物局'},
    {id:60,name:'中医药管理局'},
    {id:61,name:'疾病预防控制局'},
    {id:62,name:'矿山安全监察局'},
    {id:63,name:'消防救援局'},
    {id:64,name:'外汇管理局'},
    {id:65,name:'药品监督管理局'},
    {id:66,name:'知识产权局'},
    {id:67,name:'公务员局'},
    {id:68,name:'档案局'},
    {id:69,name:'保密局'},
    {id:70,name:'密码管理局'},
    {id:71,name:'航天局'},
    {id:72,name:'原子能机构'},
    {id:73,name:'宗教事务局'},
    {id:74,name:'台湾事务办公室'},
    {id:75,name:'乡村振兴局'},
    {id:76,name:'核安全局'},
    {id:77,name:'认证认可监督管理委员会'},
    {id:78,name:'语言文字工作委员会'},
    {id:79,name:'电影局'},
];

// 生成选项列表的函数
const generateOptions = (dataList) => dataList.map(item => (
    <option key={item.id} value={item.id}>{item.name}</option>
));

// 根据ID获取省份或部门名称的函数
const getNameById = (dataList, id) => dataList.find(item => item.id === id)?.name;

const SearchComponent = () => {
    const [query, setQuery] = useState('');
    const [pid, setPid] = useState(0); // 初始 pid 值为 0
    const [did, setDid] = useState(0); // 初始 did 值为 0
    const [results, setResults] = useState([]);
    const [currentPage, setCurrentPage] = useState(1); // 将当前页数初始值设为 1
    const [total, setTotal] = useState(0); // 总条目数
    const [summary, setSummary] = useState(''); // 添加摘要内容的状态
    const [showSummaryModal, setShowSummaryModal] = useState(false); // 添加显示摘要的状态
    const [options,setOptions] = useState([]);

    const [api, contextHolder] = notification.useNotification();
    const openNotificationWithIcon = (type) => {
        api[type]({
            message: '搜索栏不能为空',
            duration: 0.75,
            placement: 'top'
        });
    };

    useEffect(() => {
        fetchData();
    }, [currentPage]); // 当当前页数、pid 或 did 发生变化时，重新获取数据

    const handleSubmitChange = (event) => {
        setQuery(event.target.value);
    };

    const handleSearchChange = (value) => {
        setOptions(value);
        console.log(value);
    }

    const handlePidChange = (event) => {
        setPid(event.target.value);
    };

    const handleDidChange = (event) => {
        setDid(event.target.value);
    };

    const formatDate = (dateString) => {
        const date = new Date(dateString);
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        return `${year}-${month}-${day}`;
    };

    const fetchData = async () => {
        try {
            const response = await fetch(`http://localhost:3000/search?s=${query}&e=${options}&pid=${pid}&did=${did}&page=${currentPage}`);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data = await response.json();
            setResults(data.hits);
            setTotal(data.total.value)
        } catch (error) {
            console.error('Error fetching data:', error);
        }
    };

    const handleSubmit = (event) => {
        event.preventDefault();
        if (query.trim() === '') {
            // message.warning('搜索栏不能为空'); // 显示提示信息
            openNotificationWithIcon('warning')
            return;
        }
        setCurrentPage(1); // 提交搜索时重置当前页数为第一页
        fetchData(); // 在重置 currentPage 后立即执行 fetchData
    };

    const handlePageChange = (page) => {
        setCurrentPage(page); // 当分页改变时更新当前页数
    };

    const generateSummary = async (url) => {
        console.log(url)
        try {
            const response = await fetch(`http://localhost:3000/summary?url=${encodeURIComponent(url)}`);
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            const data = await response.json();
            setSummary(data.summary); // 更新摘要内容状态
            setShowSummaryModal(true); // 显示摘要弹出框
        } catch (error) {
            console.error('Error fetching summary:', error);
        }
    };

    const handleResetSummary = () => {
        setSummary(''); // 清空摘要内容状态
        setShowSummaryModal(false); // 关闭摘要弹出框
    };

    return (
        <div className="search-container">
            {contextHolder}
            <form onSubmit={handleSubmit} className="search-form">
                <input
                    type="text"
                    placeholder="Search..."
                    value={query}
                    onChange={handleSubmitChange}
                    className="search-input"
                />
                <Select
                    mode="tags"
                    style={{
                        width: '20%', // 你可以根据需要调整宽度
                    }}
                    placeholder="精确搜索词"
                    onChange={handleSearchChange}
                    value={options} // 确保 Select 的 value 与 options 状态同步
                />
                <select value={pid} onChange={handlePidChange} className="search-select">
                    {generateOptions(provincesData)}
                </select>
                <select value={did} onChange={handleDidChange} className="search-select">
                    {generateOptions(departmentsData)}
                </select>
                <button type="submit" className="search-button">Search</button>
            </form>
            <ul className="search-results">
                {results.map((result, index) => (
                    <li key={index} className="search-result">
                        <a href={result._source.url} target="_blank" rel="noopener noreferrer"
                           className="result-title">{result._source.title}</a>
                        <button onClick={() => generateSummary(result._source.url)} className="generate-summary-button">生成摘要</button>
                        <p className="result-date">Date: {formatDate(result._source.date)}</p>
                        <p className="result-content" dangerouslySetInnerHTML={{__html: result.highlight.content}}/>
                        <p className="result-source-tag">
                            {/*来源：<span className="result-province">{getNameById(provincesData, result._source.province_id)}</span> /*/}
                            来源：<span className="result-province">{getNameById(provincesData, result._source.province_id)}</span>
                            {/*<span className="result-department">{getNameById(departmentsData, result._source.department_id)}</span>*/}
                        </p>
                    </li>
                ))}
            </ul>
            <Pagination simple
                        onChange={handlePageChange}
                        current={currentPage}
                        total={total}
                        style={{marginTop: '20px', marginBottom: '20px', textAlign: 'center' }} />

            {/* 弹出框 */}
            <Modal
                title="摘要内容"
                visible={showSummaryModal}
                onCancel={handleResetSummary}
                footer={null}
            >
                <p>{summary}</p>
            </Modal>
        </div>
    );
};

export default SearchComponent;