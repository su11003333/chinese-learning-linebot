
'use client';

import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { 
collection, 
query, 
where, 
getDocs, 
doc, 
getDoc, 
setDoc, 
updateDoc,
increment
} from 'firebase/firestore';
import { db } from '@/lib/firebase';
import { publishers, grades, semesters } from '@/constants/data';
import { 
saveCharacterSearchCache, 
loadCharacterSearchCache,
saveSearchHistory,
getFormattedSearchHistory,
clearCharacterSearchCache,
clearSearchHistory
} from '@/utils/formCache';

export default function CharacterSearch() {
const { register, handleSubmit, watch, setValue, reset, formState: { errors } } = useForm();

const [searchResult, setSearchResult] = useState(null);
const [loading, setLoading] = useState(false);
const [selectedColor, setSelectedColor] = useState('pink');
const [availableLessons, setAvailableLessons] = useState([]);
const [loadingLessons, setLoadingLessons] = useState(false);
const [searchHistory, setSearchHistory] = useState([]);
const [showHistory, setShowHistory] = useState(false);
const [showCacheIndicator, setShowCacheIndicator] = useState(false);

// 監控表單值變化
const watchedPublisher = watch('publisher');
const watchedGrade = watch('grade');
const watchedSemester = watch('semester');
const watchedLesson = watch('lesson');
const watchedCharacters = watch('characters');

const colorThemes = {
  pink: {
    bg: 'bg-gradient-to-r from-pink-100 to-purple-100',
    card: 'bg-white',
    button: 'bg-gradient-to-r from-pink-400 to-purple-400 hover:from-pink-500 hover:to-purple-500',
    input: 'focus:ring-pink-300',
    title: 'bg-gradient-to-r from-pink-500 to-purple-500 bg-clip-text text-transparent',
  },
  blue: {
    bg: 'bg-gradient-to-r from-blue-100 to-green-100',
    card: 'bg-white',
    button: 'bg-gradient-to-r from-blue-400 to-green-400 hover:from-blue-500 hover:to-green-500',
    input: 'focus:ring-blue-300',
    title: 'bg-gradient-to-r from-blue-500 to-green-500 bg-clip-text text-transparent',
  },
  yellow: {
    bg: 'bg-gradient-to-r from-yellow-100 to-orange-100',
    card: 'bg-white',
    button: 'bg-gradient-to-r from-yellow-400 to-orange-400 hover:from-yellow-500 hover:to-orange-500',
    input: 'focus:ring-yellow-300',
    title: 'bg-gradient-to-r from-yellow-500 to-orange-500 bg-clip-text text-transparent',
  }
};

const theme = colorThemes[selectedColor];

// 初始化：載入快取和搜索歷史
useEffect(() => {
  const loadCachedData = () => {
    // 載入表單快取
    const cachedForm = loadCharacterSearchCache();
    if (cachedForm) {
      setValue('publisher', cachedForm.publisher || '康軒');
      setValue('grade', cachedForm.grade || 1);
      setValue('semester', cachedForm.semester || 1);
      setValue('lesson', cachedForm.lesson || 1);
      
      // 設置主題色彩
      const publisher = cachedForm.publisher || '康軒';
      if (publisher === '康軒') setSelectedColor('pink');
      else if (publisher === '南一') setSelectedColor('blue');
      else if (publisher === '翰林') setSelectedColor('yellow');
      
      setShowCacheIndicator(true);
      setTimeout(() => setShowCacheIndicator(false), 3000);
    } else {
      // 設置預設值
      setValue('publisher', '康軒');
      setValue('grade', 1);
      setValue('semester', 1);
      setValue('lesson', 1);
    }
    
    // 載入搜索歷史
    const history = getFormattedSearchHistory();
    setSearchHistory(history);
    
    // 載入初始課程列表
    loadAvailableLessons(
      cachedForm?.publisher || '康軒', 
      cachedForm?.grade || 1, 
      cachedForm?.semester || 1
    );
  };
  
  loadCachedData();
}, [setValue]);

// 監控表單變化並自動保存快取
useEffect(() => {
  const formData = {
    publisher: watchedPublisher,
    grade: watchedGrade,
    semester: watchedSemester,
    lesson: watchedLesson,
    characters: watchedCharacters
  };
  
  // 只有當表單有值時才保存
  if (watchedPublisher || watchedGrade || watchedSemester || watchedLesson || watchedCharacters) {
    saveCharacterSearchCache(formData);
  }
}, [watchedPublisher, watchedGrade, watchedSemester, watchedLesson, watchedCharacters]);

// 載入可用課程列表
const loadAvailableLessons = async (publisher, grade, semester) => {
  setLoadingLessons(true);
  try {
    const lessonsRef = collection(db, "lessons");
    const q = query(
      lessonsRef,
      where("publisher", "==", publisher),
      where("grade", "==", parseInt(grade)),
      where("semester", "==", parseInt(semester))
    );
    
    const querySnapshot = await getDocs(q);
    const lessons = [];
    
    querySnapshot.forEach((doc) => {
      const data = doc.data();
      lessons.push({
        lesson: data.lesson,
        characterCount: data.characters?.length || 0,
        id: doc.id
      });
    });
    
    lessons.sort((a, b) => a.lesson - b.lesson);
    setAvailableLessons(lessons);
  } catch (error) {
    console.error('載入課程列表失敗:', error);
    setAvailableLessons([]);
  }
  setLoadingLessons(false);
};

// 當出版社、年級、學期改變時載入可用課程
useEffect(() => {
  if (watchedPublisher && watchedGrade && watchedSemester) {
    loadAvailableLessons(watchedPublisher, watchedGrade, watchedSemester);
  }
}, [watchedPublisher, watchedGrade, watchedSemester]);

// 生成累積生字快取ID
const generateCumulativeId = (publisher, grade, semester, lesson) => {
  return `${publisher}_${grade}_${semester}_${lesson}`.toLowerCase();
};

// 生成搜索快取ID
const generateSearchId = (publisher, grade, semester, lesson, characters) => {
  const sortedChars = characters.sort().join('');
  return `${publisher}_${grade}_${semester}_${lesson}_${sortedChars}`.toLowerCase();
};

// 獲取或建立累積生字快取
const getOrCreateCumulativeCache = async (publisher, grade, semester, lesson) => {
  const cumulativeId = generateCumulativeId(publisher, grade, semester, lesson);
  
  try {
    const cacheRef = doc(db, "cumulative_characters", cumulativeId);
    const cacheDoc = await getDoc(cacheRef);
    
    if (cacheDoc.exists()) {
      const data = cacheDoc.data();
      const cacheAge = Date.now() - new Date(data.lastUpdated).getTime();
      if (cacheAge < 7 * 24 * 60 * 60 * 1000) {
        return data.characters;
      }
    }

    console.log(`重新計算 ${publisher} ${grade}年級第${semester}學期第${lesson}課的累積生字...`);
    
    const allCharacters = new Set();
    const lessonsRef = collection(db, "lessons");
    
    const q = query(lessonsRef, where("publisher", "==", publisher));
    const querySnapshot = await getDocs(q);
    
    querySnapshot.forEach((doc) => {
      const lessonData = doc.data();
      const lessonGrade = lessonData.grade;
      const lessonSemester = lessonData.semester;
      const lessonNumber = lessonData.lesson;
      
      const isInRange = (lessonGrade < grade) || 
                       (lessonGrade === grade && lessonSemester < semester) ||
                       (lessonGrade === grade && lessonSemester === semester && lessonNumber <= lesson);
      
      if (isInRange && lessonData.characters) {
        lessonData.characters.forEach(charObj => {
          if (charObj.character) {
            allCharacters.add(charObj.character);
          }
        });
      }
    });

    const charactersArray = Array.from(allCharacters);
    
    await setDoc(cacheRef, {
      publisher,
      grade,
      semester,
      lesson,
      characters: charactersArray,
      lastUpdated: new Date().toISOString(),
      totalCount: charactersArray.length,
      courseRange: `1年級上學期第1課 ~ ${grade}年級第${semester}學期第${lesson}課`,
      version: 1
    });
    
    return charactersArray;
    
  } catch (error) {
    console.error('獲取累積生字快取失敗:', error);
    throw error;
  }
};

// 獲取或建立搜索結果快取
const getOrCreateSearchCache = async (publisher, grade, semester, lesson, queryCharacters) => {
  const searchId = generateSearchId(publisher, grade, semester, lesson, queryCharacters);
  
  try {
    const searchRef = doc(db, "search_cache", searchId);
    const searchDoc = await getDoc(searchRef);
    
    if (searchDoc.exists()) {
      await updateDoc(searchRef, {
        searchCount: increment(1),
        lastSearched: new Date().toISOString()
      });
      
      return searchDoc.data().results;
    }

    console.log(`首次搜索: ${queryCharacters.join('')}，建立快取...`);
    
    const learnedCharacters = await getOrCreateCumulativeCache(publisher, grade, semester, lesson);
    
    const results = queryCharacters.map(char => ({
      character: char,
      isLearned: learnedCharacters.includes(char),
      firstAppeared: learnedCharacters.includes(char) ? "已學過" : null
    }));

    await setDoc(searchRef, {
      publisher,
      grade,
      semester,
      lesson,
      queryCharacters,
      results,
      searchCount: 1,
      lastSearched: new Date().toISOString(),
      createdAt: new Date().toISOString()
    });
    
    return results;
    
  } catch (error) {
    console.error('獲取搜索快取失敗:', error);
    throw error;
  }
};

// 從搜索歷史快速查詢
const quickSearchFromHistory = (historyItem) => {
  setValue('publisher', historyItem.publisher);
  setValue('grade', historyItem.grade);
  setValue('semester', historyItem.semester);
  setValue('lesson', historyItem.lesson);
  
  // 設置主題色彩
  if (historyItem.publisher === '康軒') setSelectedColor('pink');
  else if (historyItem.publisher === '南一') setSelectedColor('blue');
  else if (historyItem.publisher === '翰林') setSelectedColor('yellow');
  
  setShowHistory(false);
  
  // 自動提交搜索
  setTimeout(() => {
    handleSubmit(onSubmit)();
  }, 100);
};

// 主查詢函數
const onSubmit = async (data) => {
  setLoading(true);
  
  try {
    const { publisher, grade, semester, lesson, characters } = data;
    const endGrade = parseInt(grade);
    const endSemester = parseInt(semester);
    const endLesson = parseInt(lesson);
    
    const targetCharacters = characters.trim().split('').filter(char => 
      char.trim() && /[\u4e00-\u9fff]/.test(char)
    );
    
    if (targetCharacters.length === 0) {
      setSearchResult({
        error: true,
        message: '請輸入有效的中文字符'
      });
      setLoading(false);
      return;
    }

    const startTime = Date.now();
    const results = await getOrCreateSearchCache(publisher, endGrade, endSemester, endLesson, targetCharacters);
    const queryTime = Date.now() - startTime;

    const searchResult = {
      publisher,
      grade: endGrade,
      semester: endSemester,
      lesson: endLesson,
      results,
      totalLearned: results.filter(r => r.isLearned).length,
      totalQueried: results.length,
      queryTime,
      fromCache: queryTime < 500
    };

    setSearchResult(searchResult);

    // 保存到搜索歷史
    saveSearchHistory({
      publisher,
      grade: endGrade,
      semester: endSemester,
      lesson: endLesson,
      characters
    });
    
    // 更新搜索歷史顯示
    const newHistory = getFormattedSearchHistory();
    setSearchHistory(newHistory);

    // 根據出版社變更顏色主題
    if (publisher === '康軒') setSelectedColor('pink');
    else if (publisher === '南一') setSelectedColor('blue');
    else if (publisher === '翰林') setSelectedColor('yellow');
    
    setTimeout(() => {
      const resultElement = document.getElementById('search-results');
      if (resultElement) {
        resultElement.scrollIntoView({ 
          behavior: 'smooth', 
          block: 'start' 
        });
      }
    }, 300);
    
  } catch (error) {
    console.error("查詢錯誤:", error);
    setSearchResult({
      error: true,
      message: error.message || '查詢發生錯誤，請稍後再試'
    });
  }
  
  setLoading(false);
};

// 重新查詢函數
const handleResetSearch = () => {
  setSearchResult(null);
  
  setTimeout(() => {
    const formElement = document.getElementById('search-form');
    if (formElement) {
      formElement.scrollIntoView({ 
        behavior: 'smooth', 
        block: 'start' 
      });
    } else {
      window.scrollTo({ 
        top: 0, 
        behavior: 'smooth' 
      });
    }
  }, 100);
};

// 清除所有快取
const clearAllCache = () => {
  if (confirm('確定要清除所有表單快取和搜索歷史嗎？')) {
    clearCharacterSearchCache();
    clearSearchHistory();
    setSearchHistory([]);
    reset({
      publisher: '康軒',
      grade: 1,
      semester: 1,
      lesson: 1,
      characters: ''
    });
    setSelectedColor('pink');
    alert('快取已清除');
  }
};

return (
  <div className={`min-h-screen py-12 px-4 ${theme.bg}`}>
    <div className="max-w-2xl mx-auto">
      <div className="flex items-center justify-center mb-8">
        <h1 className={`text-3xl sm:text-4xl font-bold text-center ${theme.title}`}>
          累積生字查詢
        </h1>
      </div>
      
      <div id="search-form" className={`${theme.card} rounded-3xl shadow-xl p-6 mb-8 relative`}>
        {/* 操作按鈕 */}
        {searchHistory.length > 0 && (
          <div className="flex justify-start items-center mb-4">
            <button
              type="button"
              onClick={() => setShowHistory(!showHistory)}
              className="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded-full hover:bg-gray-200 transition-colors"
            >
              {showHistory ? '隱藏' : '顯示'}搜索歷史 ({searchHistory.length})
            </button>
          </div>
        )}

        {/* 搜索歷史 */}
        {showHistory && searchHistory.length > 0 && (
          <div className="mb-6 p-4 bg-gray-50 rounded-xl">
            <h3 className="text-sm font-medium text-gray-700 mb-3">🕒 最近搜索</h3>
            <div className="space-y-2 max-h-40 overflow-y-auto">
              {searchHistory.map((item, index) => (
                <button
                  key={index}
                  onClick={() => quickSearchFromHistory(item)}
                  className="w-full text-left px-3 py-2 text-sm bg-white rounded-lg hover:bg-blue-50 border border-gray-200 hover:border-blue-300 transition-colors"
                >
                  <div className="flex justify-between items-center">
                    <span className="font-medium">{item.displayText}</span>
                    <span className="text-xs text-gray-500">{item.timeAgo}</span>
                  </div>
                </button>
              ))}
            </div>
          </div>
        )}

        <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
          {/* 第一行：出版社、年級、學期 */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                出版社
              </label>
              <select
                {...register('publisher', { required: '請選擇出版社' })}
                className={`w-full px-4 py-3 rounded-full border border-gray-300 ${theme.input} focus:outline-none focus:ring-2 focus:border-transparent transition`}
                onChange={(e) => {
                  const publisher = e.target.value;
                  if (publisher === '康軒') setSelectedColor('pink');
                  else if (publisher === '南一') setSelectedColor('blue');
                  else if (publisher === '翰林') setSelectedColor('yellow');
                }}
              >
                {publishers.map(publisher => (
                  <option key={publisher} value={publisher}>{publisher}</option>
                ))}
              </select>
              {errors.publisher && (
                <p className="mt-1 text-xs text-red-500">{errors.publisher.message}</p>
              )}
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                年級
              </label>
              <select
                {...register('grade', { required: '請選擇年級' })}
                className={`w-full px-4 py-3 rounded-full border border-gray-300 ${theme.input} focus:outline-none focus:ring-2 focus:border-transparent transition`}
              >
                {grades.map(grade => (
                  <option key={grade} value={grade}>{grade}年級</option>
                ))}
              </select>
              {errors.grade && (
                <p className="mt-1 text-xs text-red-500">{errors.grade.message}</p>
              )}
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                學期
              </label>
              <select
                {...register('semester', { required: '請選擇學期' })}
                className={`w-full px-4 py-3 rounded-full border border-gray-300 ${theme.input} focus:outline-none focus:ring-2 focus:border-transparent transition`}
              >
                {semesters.map(semester => (
                  <option key={semester} value={semester}>第{semester}學期</option>
                ))}
              </select>
              {errors.semester && (
                <p className="mt-1 text-xs text-red-500">{errors.semester.message}</p>
              )}
            </div>
          </div>

          {/* 第二行：課次 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              課次 
            </label>
            <input
              type="number"
              {...register('lesson', { 
                required: '請輸入課次',
                min: { value: 1, message: '課次必須大於 0' },
                valueAsNumber: true
              })}
              className={`w-full px-4 py-3 rounded-full border border-gray-300 ${theme.input} focus:outline-none focus:ring-2 focus:border-transparent transition`}
              placeholder="第幾課"
              min="1"
            />
            {errors.lesson && (
              <p className="mt-1 text-xs text-red-500">{errors.lesson.message}</p>
            )}

          </div>
          
          {/* 第三行：要查詢的字符 */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              要查詢的字符
            </label>
            <input
              type="text"
              {...register('characters', { 
                required: '請輸入要查詢的字符',
                pattern: {
                  value: /[\u4e00-\u9fff]/,
                  message: '請輸入有效的中文字符'
                }
              })}
              className={`w-full px-4 py-3 rounded-full border border-gray-300 ${theme.input} focus:outline-none focus:ring-2 focus:border-transparent transition`}
              placeholder="例如：你我他"
            />
            {errors.characters && (
              <p className="mt-1 text-xs text-red-500">{errors.characters.message}</p>
            )}

          </div>
          
          <button
            type="submit"
            disabled={loading}
            className={`w-full py-3 px-4 ${theme.button} text-white font-bold rounded-full focus:outline-none focus:ring-2 focus:ring-opacity-50 transition disabled:opacity-70`}
          >
            {loading ? (
              <div className="flex items-center justify-center">
                <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                查詢中...
              </div>
            ) : '查詢累積生字'}
          </button>
        </form>
      </div>
      
      {/* 查詢結果 */}
      {searchResult && (
        <div id="search-results" className={`${theme.card} rounded-3xl shadow-xl p-6 mb-8 transform transition duration-300 animate-float`}>
          {searchResult.error ? (
            <div className="text-center py-8">
              <div className="flex justify-center mb-6">
                <div className="w-24 h-24 bg-red-200 rounded-full flex items-center justify-center">
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-12 w-12 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
                  </svg>
                </div>
              </div>
              <h3 className="text-xl font-bold mb-2 text-gray-800">查詢出現錯誤</h3>
              <p className="text-gray-600 mb-6">{searchResult.message}</p>
              <button
                onClick={handleResetSearch}
                className={`px-6 py-2 ${theme.button} text-white rounded-full font-medium`}
              >
                重新查詢
              </button>
            </div>
          ) : (
            <div>
              <h2 className="text-2xl font-bold mb-4 text-gray-800 text-center">查詢結果</h2>
              
              
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4 mb-6">
                {searchResult.results.map((result, index) => (
                  <div key={index} className={`
                    border-2 rounded-2xl p-4 text-center transition-all duration-200
                    ${result.isLearned 
                      ? 'border-green-300 bg-green-50' 
                      : 'border-red-300 bg-red-50'
                    }
                  `}>
                    <div className="text-4xl font-bold mb-2 text-gray-800">
                      {result.character}
                    </div>
                    <div className={`
                      inline-flex items-center px-3 py-1 rounded-full text-sm font-medium
                      ${result.isLearned
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                      }
                    `}>
                      {result.isLearned ? (
                        <>
                          <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                            <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                          </svg>
                          已學過
                        </>
                      ) : (
                        <>
                          <svg className="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                            <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd" />
                          </svg>
                          尚未學過
                        </>
                      )}
                    </div>
                  </div>
                ))}
              </div>
              
              <div className="text-center">
                <button
                  onClick={handleResetSearch}
                  className={`px-6 py-2 ${theme.button} text-white rounded-full font-medium`}
                >
                  重新查詢
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  </div>
);
}