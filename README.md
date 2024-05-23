# yamldiff

用于比较并输出两个 YAML/JSON 文件的差异的 CLI 工具。

## 使用方法
- 给定的 1.yml 文件为：
    ```yaml
    a: cat
    b: 
      - a: 1
      - b: 2
      - c: 3
    c:
      d: 1
      e: 
      f: 3
    ```

- 给定的 2.yml 文件为：
    ```yaml
    a: dog
    b: 
      - a: 1
      - b: 4
      - c: 3
    c:
      d: 1
      e: 2
      f: 4
    g: h
    ```

- 执行对比：
    ```
    yamldiff 1.yml 2.yml
    ```

- 输出结果：
    ```yaml
    a: dog
    b:
    - a: 1
    - b: 4
    - c: 3
    c:
      e: 2
      f: 4
    g: h
    ```