package controllers

import (
    "github.com/1ssk/admin-onbording/initializers"
    "github.com/1ssk/admin-onbording/models"
    "github.com/gin-gonic/gin"
)

func CreateTest(c *gin.Context) {
   
    var testData struct {
        LessonID  uint `json:"lesson_id"`
        Questions []struct {
            Text    string `json:"text"`
            Answers []struct {
                Text      string `json:"text"`
                IsCorrect bool   `json:"is_correct"`
            } `json:"answers"`
        } `json:"questions"`
    }


    if err := c.Bind(&testData); err != nil {
        c.Status(400)
        return
    }


    newTest := models.Test{
        LessonID: testData.LessonID,
    }
    result := initializers.DB.Create(&newTest)
    if result.Error != nil {
        c.Status(400)
        return
    }


    for _, q := range testData.Questions {
        newQuestion := models.Question{
            TestID: newTest.ID,
        }
        result = initializers.DB.Create(&newQuestion)
        if result.Error != nil {
            c.Status(400)
            return
        }

        for _, a := range q.Answers {
            newAnswer := models.Answer{
                QuestionID: newQuestion.ID,
                Text:       a.Text,
                IsCorrect:  a.IsCorrect,
            }
            result = initializers.DB.Create(&newAnswer)
            if result.Error != nil {
                c.Status(400)
                return
            }
        }
    }

    c.JSON(200, gin.H{
        "test": testData,
    })
}

func GetTest(c *gin.Context) {
 
    id := c.Param("id")

 
    var test models.Test
    result := initializers.DB.Preload("Questions.Answers").First(&test, id)
    if result.Error != nil {
        c.Status(404)
        return
    }

    c.JSON(200, gin.H{
        "test": test,
    })
}

func UpdateTest(c *gin.Context) {

    id := c.Param("id")

   
    var testData struct {
        LessonID  uint `json:"lesson_id"`
        Questions []struct {
            ID      uint   `json:"id"` 
            Text    string `json:"text"`
            Answers []struct {
                ID        uint   `json:"id"` 
                Text      string `json:"text"`
                IsCorrect bool   `json:"is_correct"`
            } `json:"answers"`
        } `json:"questions"`
    }


    if err := c.Bind(&testData); err != nil {
        c.Status(400)
        return
    }


    var test models.Test
    result := initializers.DB.First(&test, id)
    if result.Error != nil {
        c.Status(404)
        return
    }


    initializers.DB.Model(&test).Updates(models.Test{
        LessonID: testData.LessonID,
    })

   
    for _, q := range testData.Questions {
        if q.ID == 0 {

            newQuestion := models.Question{
                TestID: test.ID,
            }
            result = initializers.DB.Create(&newQuestion)
            if result.Error != nil {
                c.Status(400)
                return
            }

            for _, a := range q.Answers {
                newAnswer := models.Answer{
                    QuestionID: newQuestion.ID,
                    Text:       a.Text,
                    IsCorrect:  a.IsCorrect,
                }
                result = initializers.DB.Create(&newAnswer)
                if result.Error != nil {
                    c.Status(400)
                    return
                }
            }
        } else {
         
            var question models.Question
            result = initializers.DB.First(&question, q.ID)
            if result.Error != nil {
                c.Status(404)
                return
            }

            for _, a := range q.Answers {
                if a.ID == 0 {
                  
                    newAnswer := models.Answer{
                        QuestionID: question.ID,
                        Text:       a.Text,
                        IsCorrect:  a.IsCorrect,
                    }
                    result = initializers.DB.Create(&newAnswer)
                    if result.Error != nil {
                        c.Status(400)
                        return
                    }
                } else {
                
                    var answer models.Answer
                    result = initializers.DB.First(&answer, a.ID)
                    if result.Error != nil {
                        c.Status(404)
                        return
                    }
                    initializers.DB.Model(&answer).Updates(models.Answer{
                        Text:      a.Text,
                        IsCorrect: a.IsCorrect,
                    })
                }
            }
        }
    }

    c.JSON(200, gin.H{
        "test": testData,
    })
}

func DeleteTest(c *gin.Context) {
    id := c.Param("id")

    var test models.Test
    result := initializers.DB.First(&test, id)
    if result.Error != nil {
        c.JSON(404, gin.H{"error": "Данный тест не существует"})
        return
    }

   
    result = initializers.DB.Unscoped().Delete(&test)
    if result.Error != nil {
        c.JSON(400, gin.H{"error": "Не удалось удалить тест"})
        return
    }

    c.JSON(200, gin.H{
        "test": "access delete",
    })
}

